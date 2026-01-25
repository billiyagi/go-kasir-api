package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-kasir-api/internal/category"
	"go-kasir-api/internal/product"

	_ "go-kasir-api/docs" // Import generated docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Kasir API
// @version 1.0
// @description This is a simple cashier API service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	// Buat multiplexer (router) baru.
	// Ini memungkinkan kita mendaftarkan handler secara modular.
	mux := http.NewServeMux()

	// Daftarkan handler untuk Swagger UI.
	// URL-nya akan menjadi http://localhost:8080/swagger/index.html
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Daftarkan health check endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Daftarkan semua endpoint dari package product.
	product.RegisterHandlers(mux)

	// Daftarkan semua endpoint dari package category.
	category.RegisterHandlers(mux)

	// Gunakan mux yang sudah dikonfigurasi untuk server.
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server listening on port 8080")
	log.Fatal(server.ListenAndServe())
}