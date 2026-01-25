package category

import (
	"encoding/json"
	"net/http"
)

// Category represents a product category.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RegisterHandlers sets up the routing for the category endpoints.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/categories", categoriesHandler)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category endpoint is not implemented yet.",
	})
}
