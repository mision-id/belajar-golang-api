package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/models"
	"kasir-api/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// Handler methods GET /api/categories, POST /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getAll(w, r)
	case "POST":
		h.CreateCategory(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// function to show all categories
func (h *CategoryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, "Failed to encode categories", http.StatusInternalServerError)
		return
	}
}

// function to create a new category
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateCategory(&category)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler methods GET /api/categories/{id}, PUT /api/categories/{id}, DELETE /api/categories/{id}
func (h *CategoryHandler) HandleCategoriesByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		h.GetByID(w, r, id)
	case "PUT":
		h.UpdateByID(w, r, id)
	case "DELETE":
		h.DeleteByID(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// function to get category by id
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve category", http.StatusInternalServerError)
		return
	}
	if category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		http.Error(w, "Failed to encode category", http.StatusInternalServerError)
	}
}

// function to update category by id
func (h *CategoryHandler) UpdateByID(w http.ResponseWriter, r *http.Request, id int) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateByID(id, &category)
	if err != nil {
		if err.Error() == "no rows updated" {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// function to delete category by id
func (h *CategoryHandler) DeleteByID(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteByID(id)
	if err != nil {
		if err.Error() == "Category not found" {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
	if err != nil {
		http.Error(w, "Failed to encode delete response", http.StatusInternalServerError)
	}
}
