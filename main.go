package main

// @title Kasir API
// @version 1.0
// @description Simple RESTful API for managing products and categories.
// @host cwu-go-kasirapi.zeabur.app
// @BasePath /
// @schemes https http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "kasir-api/docs"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

var categories = []Category{
	{ID: 1, Name: "Food", Description: "Food products"},
	{ID: 2, Name: "Drink", Description: "Drink products"},
}

var products = []Product{
	{ID: 1, Name: "Indomie Jumbo", Price: 100, Stock: 10, CategoryID: 1},
	{ID: 2, Name: "Vitamin C 100ml", Price: 200, Stock: 20, CategoryID: 1},
	{ID: 3, Name: "Sambal ABC", Price: 300, Stock: 30, CategoryID: 1},
	{ID: 4, Name: "Coca Cola", Price: 400, Stock: 40, CategoryID: 2},
	{ID: 5, Name: "Sprite", Price: 500, Stock: 50, CategoryID: 2},
	{ID: 6, Name: "Pepsi", Price: 600, Stock: 60, CategoryID: 2},
}

// getCategoryByID godoc
// @Summary Get a category
// @Description Get a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 404 {string} string "Category not found"
// @Router /api/categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL path
	// URL: /api/categories/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Find category by ID
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// getCategories godoc
// @Summary List categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /api/categories [get]
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// createCategory godoc
// @Summary Create a category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "New Category"
// @Success 201 {object} Category
// @Router /api/categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
	// read data from request body
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// add new category to categories
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newCategory)
}

// updateCategory godoc
// @Summary Update a category
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Updated Category"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Category not found"
// @Router /api/categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id from request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// get data from request body
	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// update category
	for i, c := range categories {
		if c.ID == id {
			updatedCategory.ID = id
			categories[i] = updatedCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// deleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "Category not found"
// @Router /api/categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id from request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// delete category
	for i, c := range categories {
		if c.ID == id {
			// remove category from slice by creating new slice without the deleted category
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// getProducts godoc
// @Summary List products
// @Description Get all products
// @Tags products
// @Produce json
// @Success 200 {array} Product
// @Router /api/products [get]
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// createProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "New Product"
// @Success 201 {object} Product
// @Router /api/products [post]
func createProduct(w http.ResponseWriter, r *http.Request) {
	// read data from request body
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// add new product to products
	newProduct.ID = len(products) + 1
	products = append(products, newProduct)

	// return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newProduct)
}

// getProductByID godoc
// @Summary Get a product
// @Description Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 404 {string} string "Product not found"
// @Router /api/products/{id} [get]
func getProductByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL path
	// URL: /api/products/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// Find product by ID
	for _, p := range products {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// updateProduct godoc
// @Summary Update a product
// @Description Update a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product true "Updated Product"
// @Success 200 {object} Product
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Product not found"
// @Router /api/products/{id} [put]
func updateProduct(w http.ResponseWriter, r *http.Request) {
	// get id from request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// get data from request body
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// update product
	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// deleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {string} string "Product not found"
// @Router /api/products/{id} [delete]
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	// get id from request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	// delete product
	for i, p := range products {
		if p.ID == id {
			// remove product from slice by creating new slice without the deleted product
			products = append(products[:i], products[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func main() {
	// GET /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
	})

	// GET /api/products
	// POST /api/products
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProducts(w, r)
		case "POST":
			createProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/products/{id}
	// PUT /api/products/{id}
	// DELETE /api/products/{id}
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getProductByID(w, r)
		case "PUT":
			updateProduct(w, r)
		case "DELETE":
			deleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/categories
	// POST /api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategories(w, r)
		case "POST":
			createCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/categories/{id}
	// PUT /api/categories/{id}
	// DELETE /api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCategoryByID(w, r)
		case "PUT":
			updateCategory(w, r)
		case "DELETE":
			deleteCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve Swagger JSON
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Serve Scalar API Reference
	http.HandleFunc("/reference", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/swagger.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		w.Write([]byte(html))
	})

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on :%s\n", port)

	error := http.ListenAndServe(":"+port, nil)
	if error != nil {
		fmt.Println("Error starting server: ", error)
	}
}
