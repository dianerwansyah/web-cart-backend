package setup

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

	log.Printf("Starting Setup service on port %s...", cfg.Server.SetupPort)
	if err := http.ListenAndServe(":"+cfg.Server.SetupPort, corsRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
