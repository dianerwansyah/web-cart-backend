package logic

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProducts() ([]model.Product, error) {
	collectionName := helper.GetTableName(model.Product{})
	collection := helper.GetCollection(collectionName)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return nil, err
	}

	var products []model.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return nil, err
	}
	return products, nil
}

type FilterRequest struct {
	FilterType string `json:"filtertype"`
	Limit      int    `json:"limit"`
	Field      string `json:"field,omitempty"`
	Value      string `json:"value,omitempty"`
}

func GetProductsByFilter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filterRequest FilterRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&filterRequest); err != nil {
		log.Printf("Error decoding request body: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	collectionName := helper.GetTableName(model.Product{})
	collection := helper.GetCollection(collectionName)

	var cursor *mongo.Cursor
	var err error

	if filterRequest.FilterType == "limit" {
		// Find products with limit if limit is greater than 0
		if filterRequest.Limit > 0 {
			options := options.Find().SetLimit(int64(filterRequest.Limit))
			cursor, err = collection.Find(ctx, bson.M{}, options)
		} else {
			// Find all products without limit
			cursor, err = collection.Find(ctx, bson.M{})
		}
	} else {
		// Filter by specific field and value
		filter := bson.M{filterRequest.Field: filterRequest.Value}
		cursor, err = collection.Find(ctx, filter)
	}

	if err != nil {
		log.Printf("Error finding products: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error finding products")
		return
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error decoding products")
		return
	}

	// Pilih produk secara acak hingga mencapai limit jika limit > 0
	if filterRequest.FilterType == "limit" && filterRequest.Limit > 0 && len(products) > filterRequest.Limit {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(products), func(i, j int) { products[i], products[j] = products[j], products[i] })
		products = products[:filterRequest.Limit]
	}

	helper.RespondWithJSON(w, http.StatusOK, products)
}

type UserRequest struct {
	UserID     string `json:"UserID"`
	IsCheckout bool   `json:"IsCheckout"`
	IsConfirm  bool   `json:"IsConfirm"`
}

type ProductWithQuantity struct {
	Product  model.Product `json:"Product"`
	Quantity int           `json:"Quantity"`
}

func GetProductsUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userRequest UserRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		log.Printf("Error decoding request body: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user ID from the request
	userID, err := primitive.ObjectIDFromHex(userRequest.UserID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Create filter for transactions
	filter := bson.M{"UserID": userID, "IsCheckout": userRequest.IsCheckout}

	// Fetch transactions for the user
	transactionCollection := helper.GetCollection(helper.GetTableName(model.Cart{}))
	var transactions []model.Cart
	cursor, err := transactionCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding transactions: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error finding transactions")
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &transactions); err != nil {
		log.Printf("Error decoding transactions: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error decoding transactions")
		return
	}

	// Create a map to store product quantities
	productQuantities := make(map[primitive.ObjectID]int)
	for _, transaction := range transactions {
		productQuantities[transaction.ProductID] += transaction.Quantity
	}

	// Extract product IDs from the transactions
	var productIDs []primitive.ObjectID
	for productID := range productQuantities {
		productIDs = append(productIDs, productID)
	}

	// Fetch products based on the extracted product IDs
	productCollection := helper.GetCollection(helper.GetTableName(model.Product{}))
	var products []model.Product
	cursor, err = productCollection.Find(ctx, bson.M{"_id": bson.M{"$in": productIDs}})
	if err != nil {
		log.Printf("Error finding products: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error finding products")
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error decoding products")
		return
	}

	// Combine products with their quantities
	var productsWithQuantity []ProductWithQuantity
	for _, product := range products {
		productsWithQuantity = append(productsWithQuantity, ProductWithQuantity{
			Product:  product,
			Quantity: productQuantities[product.ID],
		})
	}

	helper.RespondWithJSON(w, http.StatusOK, productsWithQuantity)
}

type UpdateCartItemRequest struct {
	ProductID string `json:"ProductID"`
	Quantity  int    `json:"Quantity"`
	UserID    string `json:"UserID"`
}

func UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	var req UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	collection := helper.GetCollection(model.Cart{}.TableName())

	// Check if the transaction already exists
	filter := bson.M{"ProductID": productID, "UserID": userID}
	var existingTransaction model.Cart
	err = collection.FindOne(context.Background(), filter).Decode(&existingTransaction)
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Error checking existing transaction", http.StatusInternalServerError)
		return
	}

	// If the transaction exists, update it
	if existingTransaction.ID != primitive.NilObjectID {
		if req.Quantity > 0 {
			update := bson.M{
				"$set": bson.M{
					"Quantity": req.Quantity,
					"Created":  time.Now(),
				},
			}
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				http.Error(w, "Error updating transaction", http.StatusInternalServerError)
				return
			}
			existingTransaction.Quantity = req.Quantity
			existingTransaction.Created = time.Now()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingTransaction)
			return
		} else {
			_, err := collection.DeleteOne(context.Background(), filter)
			if err != nil {
				http.Error(w, "Error deleting transaction", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
			return
		}
	}

	// If the transaction does not exist, insert a new one
	transaction := model.Cart{
		ProductID: productID,
		UserID:    userID,
		Quantity:  req.Quantity,
		Created:   time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), transaction)
	if err != nil {
		http.Error(w, "Error saving transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
