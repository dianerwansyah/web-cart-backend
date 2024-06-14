package setup

import (
	"github.com/dianerwansyah/web-cart-backend/helper"
	"github.com/dianerwansyah/web-cart-backend/logic"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/products/gets", helper.GenericGetHandler(helper.ConvertToInterface(logic.GetProducts))).Methods("GET")
	r.HandleFunc("/api/products/get", logic.GetProductsByFilter).Methods("POST")
	r.HandleFunc("/api/categories", helper.GenericGetHandler(helper.ConvertToInterface(logic.GetCategories))).Methods("GET")
	r.HandleFunc("/api/categories/{id}", logic.GetCategoryByID).Methods("GET")
	r.HandleFunc("/api/cart/save", logic.UpdateCartItemQuantity).Methods("POST")
	r.HandleFunc("/api/cart/get", logic.GetProductsUser).Methods("POST")
	r.HandleFunc("/api/cart/savecheckout", logic.SaveCheckout).Methods("POST")
	r.HandleFunc("/api/cart/saveconfirm", logic.SaveConfirm).Methods("POST")
	r.HandleFunc("/api/history/get", logic.GetHistory).Methods("POST")
	return r
}
