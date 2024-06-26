package iam

import (
	"log"
	"net/http"

	"github.com/dianerwansyah/web-cart-backend/helper"
)

func StartServer() {
	cfg := helper.GetConfig()
	r := NewRouter()

	// Apply CORS middleware
	corsRouter := helper.CORSMiddleware(r)

	log.Printf("Starting IAM service on port %s...", cfg.Server.IamPort)
	if err := http.ListenAndServe(":"+cfg.Server.IamPort, corsRouter); err != nil {
		log.Fatalf("Failed to start IAM server: %v", err)
	}
}
