package iam

import (
	"github.com/dianerwansyah/web-cart/backend/logic"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", logic.LoginHandler).Methods("POST")
	return r
}
