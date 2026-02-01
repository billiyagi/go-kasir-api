package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go-kasir-api/database"
	"go-kasir-api/handlers"
	"go-kasir-api/repositories"
	"go-kasir-api/services"

	_ "go-kasir-api/docs" // Import generated docs

	viper "github.com/spf13/viper"
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
	// 1. Initialize Viper and Load Config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading config file: %v", err)
		}
	}

	type Config struct {
		Port   string `mapstructure:"PORT"`
		DBConn string `mapstructure:"DB_CONN"`
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Default port if not set
	if config.Port == "" {
		config.Port = "8080"
	}

	// 2. Setup Database
	if config.DBConn == "" {
		log.Fatal("DB_CONN is not set in environment or config")
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 3. Initialize Injectors (Repositories, Services, Handlers)
	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// 4. Setup Router
	mux := http.NewServeMux()

	// 5. Register Routes
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Product Routes
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Category Routes
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Package specific routes (Legacy - can be removed if fully migrated)
	// product.RegisterHandlers(mux) // Legacy removed
	// category.RegisterHandlers(mux) // Removed legacy category handler

	// 6. Start Server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	fmt.Printf("Server listening on port %s\n", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}