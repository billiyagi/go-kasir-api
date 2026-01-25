package product

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Product represents a product in the cashier system
type Product struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// In-memory storage
var products = []Product{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}

// nextID helps in assigning a new ID for new products.
var nextID = 4

// RegisterHandlers sets up the routing for the product endpoints.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/products", productsHandler)
	mux.HandleFunc("/api/products/", productHandler)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProductByID(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getProducts godoc
// @Summary Get all products
// @Description Get a list of all products in inventory
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /products [get]
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// createProduct godoc
// @Summary Create a new product
// @Description Add a new product to the inventory
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product object that needs to be added"
// @Success 201 {object} Product
// @Failure 400 {object} map[string]string
// @Router /products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}

// getProductByID godoc
// @Summary Get a product by ID
// @Description Get details of a single product by its integer ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Product ID"}`, http.StatusBadRequest)
		return
	}

	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, `{"error": "Product not found"}`, http.StatusNotFound)
}

// updateProduct godoc
// @Summary Update an existing product
// @Description Update details of an existing product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Product object with updated data"
// @Success 200 {object} Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Product ID"}`, http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	for i := range products {
		if products[i].ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}

	http.Error(w, `{"error": "Product not found"}`, http.StatusNotFound)
}

// deleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a single product by its integer ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Product ID"}`, http.StatusBadRequest)
		return
	}

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, `{"error": "Product not found"}`, http.StatusNotFound)
}
