package main

import (
	"context"
	"log"

	"github.com/dianerwansyah/web-cart-backend/app/iam"
	"github.com/dianerwansyah/web-cart-backend/app/setup"
	"github.com/dianerwansyah/web-cart-backend/helper"
)

func main() {
	// Inisialisasi klien MongoDB
	models := helper.GetAllModels()
	tableNames := helper.GetTableNames(models)

	cfg := helper.GetConfig()
	client := helper.InitDB(cfg.Server.MongoURI, tableNames)

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Simpan klien database ke helper untuk digunakan di seluruh aplikasi
	helper.SetDBClient(client)

	go iam.StartServer()
	go setup.StartServer()

	// Prevent the main function from exiting
	select {}
}
