package logic

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/model"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client
var userCollection *mongo.Collection

func init() {
	cfg := helper.GetConfig()
	clientOptions := options.Client().ApplyURI(cfg.Server.MongoURI)
	client, _ = mongo.Connect(context.Background(), clientOptions)
	userCollection = client.Database(cfg.Server.MongoDB).Collection("users")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser model.User
	err = userCollection.FindOne(context.Background(), bson.M{"username": creds.Username}).Decode(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := model.User{
		ID:       primitive.NewObjectID(),
		Username: creds.Username,
		Password: string(hashedPassword),
		Role:     creds.Role,
		Created:  time.Now(),
	}

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user model.User
	err = userCollection.FindOne(context.Background(), bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		log.Printf("User not found: %s", creds.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Tambahkan log untuk debugging
	log.Printf("Login attempt for user: %s", creds.Username)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		// Tambahkan log untuk debugging
		log.Printf("Password mismatch for user: %s", creds.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(helper.GetConfig().Server.JwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token and user ID
	response := map[string]string{
		"token":    tokenString,
		"userId":   user.ID.Hex(),
		"username": user.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
