package logic

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HistoryWithProduct struct {
	model.History `bson:",inline"`
	Product       model.Product `bson:"Product" json:"Product"`
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request struct {
		UserID string `json:"userId"`
	}

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer client.Disconnect(ctx)

	historyCollection := helper.GetCollection(model.History{}.TableName())
	productCollection := helper.GetCollection(model.Product{}.TableName())

	filter := bson.M{"UserID": userID}
	cursor, err := historyCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error finding history items", http.StatusInternalServerError)
		return
	}

	var historyItems []model.History
	if err = cursor.All(ctx, &historyItems); err != nil {
		http.Error(w, "Error decoding history items", http.StatusInternalServerError)
		return
	}

	var result []HistoryWithProduct
	for _, item := range historyItems {
		var product model.Product
		err := productCollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			http.Error(w, "Error finding product", http.StatusInternalServerError)
			return
		}
		result = append(result, HistoryWithProduct{
			History: item,
			Product: product,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
