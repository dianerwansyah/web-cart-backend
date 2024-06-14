package logic

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductQuantity represents a product with its quantity in the checkout payload
type ProductQuantity struct {
	ProductID   primitive.ObjectID `json:"ProductID" bson:"ProductID"`
	Name        string             `json:"Name" bson:"Name"`
	Description string             `json:"Description" bson:"Description"`
	Price       float64            `json:"Price" bson:"Price"`
	ImageURL    string             `json:"ImageURL" bson:"ImageURL"`
	Quantity    int                `json:"Quantity" bson:"Quantity"`
}

// CheckoutRequest represents the entire checkout payload
type CheckoutRequest struct {
	UserID       string            `json:"UserID" bson:"UserID"`
	IsCheckout   bool              `json:"IsCheckout" bson:"IsCheckout"`
	IsConfirm    bool              `json:"IsConfirm"`
	Target       []ProductQuantity `json:"Target" bson:"Target"`
	TotalCoupons int               `json:"TotalCoupons" bson:"TotalCoupons"`
}

func SaveCheckout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var checkoutRequest CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&checkoutRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(checkoutRequest.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cartCollection := helper.GetCollection(model.Cart{}.TableName())
	productCollection := helper.GetCollection(model.Product{}.TableName())

	for _, item := range checkoutRequest.Target {
		productID, err := primitive.ObjectIDFromHex(item.ProductID.Hex())
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Update or insert into cart
		filter := bson.M{"ProductID": productID, "UserID": userID}
		update := bson.M{
			"$set": bson.M{
				"Quantity":   item.Quantity,
				"IsCheckout": true,
				"IsConfirm":  false,
				"Created":    time.Now(),
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err = cartCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			http.Error(w, "Error updating cart", http.StatusInternalServerError)
			return
		}

		// Update product stock
		var product model.Product
		err = productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		newStock := product.Stock - item.Quantity
		if newStock < 0 {
			http.Error(w, "Insufficient stock", http.StatusBadRequest)
			return
		}

		_, err = productCollection.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": bson.M{"Stock": newStock}})
		if err != nil {
			http.Error(w, "Error updating product stock", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Checkout successful"})
}

func SaveConfirm(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var checkoutRequest CheckoutRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&checkoutRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(checkoutRequest.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cartCollection := helper.GetCollection(model.Cart{}.TableName())
	couponCollection := helper.GetCollection(model.Coupon{}.TableName())
	historyCollection := helper.GetCollection(model.History{}.TableName())

	for _, item := range checkoutRequest.Target {
		productID := item.ProductID

		// Update cart to set IsConfirm to true
		filter := bson.M{"ProductID": productID, "UserID": userID}
		update := bson.M{
			"$set": bson.M{
				"Quantity":   item.Quantity,
				"IsCheckout": true,
				"IsConfirm":  true,
				"Created":    time.Now(),
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err = cartCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			http.Error(w, "Error updating cart", http.StatusInternalServerError)
			return
		}
	}

	// Check if the coupon already exists
	var existingCoupon model.Coupon
	err = couponCollection.FindOne(ctx, bson.M{"UserID": userID}).Decode(&existingCoupon)
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Error checking existing coupon", http.StatusInternalServerError)
		return
	}

	if existingCoupon.ID != primitive.NilObjectID {
		// Update existing coupon
		newAmount := existingCoupon.Amount + checkoutRequest.TotalCoupons
		update := bson.M{
			"$set": bson.M{
				"Amount":      newAmount,
				"LastUpdated": time.Now(),
			},
		}
		_, err = couponCollection.UpdateOne(ctx, bson.M{"_id": existingCoupon.ID}, update)
		if err != nil {
			http.Error(w, "Error updating coupon", http.StatusInternalServerError)
			return
		}
	} else {
		// Insert new coupon
		coupon := model.Coupon{
			UserID:      userID,
			Amount:      checkoutRequest.TotalCoupons,
			Created:     time.Now(),
			LastUpdated: time.Now(),
		}
		_, err = couponCollection.InsertOne(ctx, coupon)
		if err != nil {
			http.Error(w, "Error inserting coupon", http.StatusInternalServerError)
			return
		}
	}

	// Move confirmed items from cart to history
	filter := bson.M{"UserID": userID, "IsCheckout": true, "IsConfirm": true}
	cursor, err := cartCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error finding confirmed cart items", http.StatusInternalServerError)
		return
	}
	var confirmedItems []model.Cart
	if err = cursor.All(ctx, &confirmedItems); err != nil {
		http.Error(w, "Error decoding confirmed cart items", http.StatusInternalServerError)
		return
	}

	// Generate a new transaction ID
	idTrx := primitive.NewObjectID()

	var historyItems []interface{}
	for _, item := range confirmedItems {
		history := model.History{
			IDTrx:      idTrx,
			ProductID:  item.ProductID,
			UserID:     item.UserID,
			Quantity:   item.Quantity,
			IsCheckout: item.IsCheckout,
			IsConfirm:  item.IsConfirm,
			Created:    item.Created,
		}
		historyItems = append(historyItems, history)
	}

	if len(historyItems) > 0 {
		_, err = historyCollection.InsertMany(ctx, historyItems)
		if err != nil {
			http.Error(w, "Error inserting history items", http.StatusInternalServerError)
			return
		}
	}

	// Delete confirmed items from cart
	_, err = cartCollection.DeleteMany(ctx, filter)
	if err != nil {
		http.Error(w, "Error deleting confirmed cart items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Confirmation successful"})
}
