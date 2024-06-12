package iam

import (
	"log"
	"net/http"

	"github.com/dianerwansyah/web-cart/backend/helper"
)

func StartServer() {
	cfg := helper.GetConfig()
	r := NewRouter()
	log.Printf("Starting IAM service on port %s...", cfg.Server.IamPort)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.IamPort, r))
}
