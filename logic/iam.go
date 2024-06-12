package logic

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dianerwansyah/web-cart/backend/helper"
	"github.com/dianerwansyah/web-cart/backend/model"
	"github.com/golang-jwt/jwt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Authenticate the user (this is just a placeholder, implement your own logic)
	if creds.Username != "user" || creds.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	cfg := helper.GetConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.Server.JwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}
