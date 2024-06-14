package logic

import (
	"context"
	"log"
	"net/http"

	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCategories() ([]model.Category, error) {
	// Mengambil nama koleksi secara dinamis dari model Category
	collectionName := helper.GetTableName(model.Category{})
	collection := helper.GetCollection(collectionName)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return nil, err
	}

	var categories []model.Category
	if err := cursor.All(context.Background(), &categories); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return nil, err
	}

	return categories, nil
}

func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	categoryID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	collection := helper.GetCollection("categories")
	filter := bson.M{"_id": categoryID}

	var category model.Category
	err = collection.FindOne(context.Background(), filter).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helper.RespondWithError(w, http.StatusNotFound, "Category not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, category)
}
