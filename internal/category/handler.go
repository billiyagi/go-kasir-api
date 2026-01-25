package category

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Category represents a product category.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// In-memory storage for categories
var categories = []Category{
	{ID: 1, Name: "Makanan"},
	{ID: 2, Name: "Minuman"},
}
var nextID = 3

// RegisterHandlers sets up the routing for the category endpoints.
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/categories", categoriesHandler)
	mux.HandleFunc("/api/categories/", categoryHandler)
}

// categoriesHandler handles requests for the categories collection
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCategories(w, r)
	case http.MethodPost:
		createCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// categoryHandler handles requests for a single category
func categoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCategoryByID(w, r)
	case http.MethodPut:
		updateCategory(w, r)
	case http.MethodDelete:
		deleteCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getCategories godoc
// @Summary Get all categories
// @Description Get a list of all product categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// createCategory godoc
// @Summary Create a new category
// @Description Add a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "Category object that needs to be added"
// @Success 201 {object} Category
// @Failure 400 {object} map[string]string
// @Router /categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	newCategory.ID = nextID
	nextID++
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// deleteCategory godoc
// @Summary Delete a category by ID
// @Description Delete a single category by its integer ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Category ID"}`, http.StatusBadRequest)
		return
	}

	for i, p := range categories {
		if p.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
}

// updateCategory godoc
// @Summary Update an existing category
// @Description Update details of an existing category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Category object with updated data"
// @Success 200 {object} Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Category ID"}`, http.StatusBadRequest)
		return
	}

	var updatedCategory Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
}

// getCategoryByID godoc
// @Summary Get a category by ID
// @Description Get details of a single category by its integer ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid Category ID"}`, http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, `{"error": "Category not found"}`, http.StatusNotFound)
}