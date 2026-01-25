package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-kasir-api/internal/category"
	"go-kasir-api/internal/product"
)

func main() {
	// Buat multiplexer (router) baru.
	// Ini memungkinkan kita mendaftarkan handler secara modular.
	mux := http.NewServeMux()

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