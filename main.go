package main

// @title Kasir API
// @version 1.0
// @description Simple RESTful API for managing products and categories.
// @host cwu-go-kasirapi-indraarianggi1211-nq45hljg.leapcell.dev
// @BasePath /
// @schemes https

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

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

	// initialize category repositories, services, and handlers
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Category routes
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// GET /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
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
