package main

import (
	"github.com/dianerwansyah/web-cart/backend/iam"
)

func main() {
	go iam.StartServer()
	// go trx.StartServer()
	// go setup.StartServer()

	// Prevent the main function from exiting
	select {}
}
