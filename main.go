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
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Food", Description: "Food products"},
	{ID: 2, Name: "Drink", Description: "Drink products"},
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

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// load config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// initialize database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// initialize product repositories, services, and handlers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// GET /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
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

	fmt.Printf("Starting server on :%s\n", config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
