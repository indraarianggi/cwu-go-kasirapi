# Kasir API

Kasir API is a simple RESTful API built with Go (Golang) for managing products and categories. This project demonstrates basic CRUD (Create, Read, Update, Delete) operations using the standard `net/http` package without any external frameworks.

## Features

- **Product Management**: Create, read, update, and delete products.
- **Category Management**: Create, read, update, and delete product categories.
- **In-Memory Storage**: Data is stored in memory (slices) for simplicity and demonstration purposes.
- **RESTful Endpoints**: Follows standard REST conventions.
- **JSON Support**: All communication is done via JSON.

## Prerequisites

- Go (Golang) installed on your machine (version 1.25.6 or later recommended).

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd kasir-api
   ```

2. Run the application:

   ```bash
   go run main.go
   ```

   The server will start on port `8080`.

## API Endpoints

### Health Check

- **GET** `/health`: Check if the API is running.

### Products

- **GET** `/api/products`: Retrieve a list of all products.
- **POST** `/api/products`: Create a new product.
  - Body:
    ```json
    {
      "name": "Product Name",
      "price": 1000,
      "stock": 10,
      "category_id": 1
    }
    ```
- **GET** `/api/products/{id}`: Retrieve a specific product by ID.
- **PUT** `/api/products/{id}`: Update an existing product.
  - Body:
    ```json
    {
      "name": "Updated Name",
      "price": 1500,
      "stock": 5,
      "category_id": 1
    }
    ```
- **DELETE** `/api/products/{id}`: Delete a product.

### Categories

- **GET** `/api/categories`: Retrieve a list of all categories.
- **POST** `/api/categories`: Create a new category.
  - Body:
    ```json
    {
      "name": "Category Name",
      "description": "Category Description"
    }
    ```
- **GET** `/api/categories/{id}`: Retrieve a specific category by ID.
- **PUT** `/api/categories/{id}`: Update an existing category.
  - Body:
    ```json
    {
      "name": "Updated Category",
      "description": "Updated Description"
    }
    ```
- **DELETE** `/api/categories/{id}`: Delete a category.

## API Documentation

This project uses OpenAPI (Swagger) 2.0 to document its endpoints and Scalar to provide a modern documentation UI.

### Viewing Documentation

Once the server is running, you can access the interactive API reference at:

- **Scalar UI**: [http://localhost:8080/reference](http://localhost:8080/reference) (Recommended)
- **Swagger JSON**: [http://localhost:8080/docs/swagger.json](http://localhost:8080/docs/swagger.json)

### Generating Documentation

If you modify the API endpoints or the comments in `main.go`, you must regenerate the documentation:

1.  **Install Swag CLI** (if not already installed):
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2.  **Generate Docs**:
    Run the following command in the project root:
    ```bash
    swag init
    ```
    This will update the files in the `docs/` directory.

## Project Structure

- `main.go`: Contains the main application logic, data structures, and route handlers.
- `go.mod`: Go module definition file.
