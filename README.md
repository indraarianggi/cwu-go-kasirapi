# Kasir API

Kasir API is a RESTful API built with Go (Golang) for managing a point-of-sale (POS) system. This project demonstrates clean architecture principles with a layered structure, implementing CRUD operations, transaction management, and sales reporting using PostgreSQL database.

## Features

- **Product Management**: Full CRUD operations for products with category relationships and name-based filtering.
- **Category Management**: Full CRUD operations for product categories.
- **Transaction System**: Complete checkout functionality with transaction tracking and detail records.
- **Sales Reporting**: Generate sales summary reports including:
  - Total revenue
  - Transaction count
  - Best-selling product
  - Optional date range filtering
- **Database Storage**: PostgreSQL database with connection pooling using pgx driver.
- **Clean Architecture**: Follows layered architecture pattern (Models → Repositories → Services → Handlers).
- **Configuration Management**: Environment-based configuration using Viper.
- **RESTful Endpoints**: Follows standard REST conventions.
- **API Documentation**: OpenAPI (Swagger) 2.0 specification with Scalar UI.
- **JSON Support**: All communication is done via JSON.

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
kasir-api/
├── main.go              # Application entry point and routing
├── database/            # Database initialization and connection
├── models/              # Data structures and domain models
├── repositories/        # Data access layer (database operations)
├── services/            # Business logic layer
├── handlers/            # HTTP handlers and request/response handling
└── docs/                # Auto-generated Swagger documentation
```

## Prerequisites

- Go (Golang) version 1.25.6 or later
- PostgreSQL database
- Swag CLI tool (for generating API documentation)

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd kasir-api
   ```

2. Set up environment variables by creating a `.env` file:

   ```bash
   PORT=8080
   DB_CONN=postgres://username:password@localhost:5432/kasir_db?sslmode=disable
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

   The server will start on the configured port (default: `8080`).

## API Endpoints

### Health Check

- **GET** `/health`: Check if the API is running.
  - Response: `{"status": "OK", "message": "API Running"}`

### Products

- **GET** `/api/products`: Retrieve a list of all products with optional name filtering.
  - Query Parameters:
    - `name` (optional): Filter products by name (partial match)
  - Example: `/api/products?name=milk`
- **POST** `/api/products`: Create a new product.
  - Request Body:
    ```json
    {
      "name": "Product Name",
      "price": 1000,
      "stock": 10,
      "category_id": 1
    }
    ```
- **GET** `/api/products/{id}`: Retrieve a specific product by ID (includes category information).
- **PUT** `/api/products/{id}`: Update an existing product.
  - Request Body:
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
  - Request Body:
    ```json
    {
      "name": "Category Name",
      "description": "Category Description"
    }
    ```
- **GET** `/api/categories/{id}`: Retrieve a specific category by ID.
- **PUT** `/api/categories/{id}`: Update an existing category.
  - Request Body:
    ```json
    {
      "name": "Updated Category",
      "description": "Updated Description"
    }
    ```
- **DELETE** `/api/categories/{id}`: Delete a category.

### Transactions

- **POST** `/api/checkout`: Process a checkout transaction with multiple items.
  - Request Body:
    ```json
    {
      "items": [
        {
          "product_id": 1,
          "quantity": 2
        },
        {
          "product_id": 3,
          "quantity": 1
        }
      ]
    }
    ```
  - Response: Returns transaction details including total amount and individual item subtotals.
  - Note: Automatically reduces product stock and creates transaction records.

### Reports

- **GET** `/api/report`: Retrieve sales summary report.
  - Query Parameters:
    - `start_date` (optional): Start date in `YYYY-MM-DD` format (e.g., `2026-01-01`)
    - `end_date` (optional): End date in `YYYY-MM-DD` format (e.g., `2026-01-31`)
  - Example: `/api/report?start_date=2026-01-01&end_date=2026-01-31`
  - Response:
    ```json
    {
      "total_revenue": 50000,
      "total_transaction": 15,
      "best_selling_product": {
        "product_id": 1,
        "product_name": "Product Name",
        "count": 25
      }
    }
    ```
  - Note: If no dates are provided, returns all-time statistics.

## Database Setup

1. Create a PostgreSQL database:
   ```sql
   CREATE DATABASE kasir_db;
   ```

2. Create the required tables:
   ```sql
   -- Categories table
   CREATE TABLE categories (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       description TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Products table
   CREATE TABLE products (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       price INTEGER NOT NULL,
       stock INTEGER NOT NULL DEFAULT 0,
       category_id INTEGER REFERENCES categories(id),
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Transactions table
   CREATE TABLE transactions (
       id SERIAL PRIMARY KEY,
       total_amount INTEGER NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Transaction details table
   CREATE TABLE transaction_details (
       id SERIAL PRIMARY KEY,
       transaction_id INTEGER REFERENCES transactions(id),
       product_id INTEGER REFERENCES products(id),
       quantity INTEGER NOT NULL,
       subtotal INTEGER NOT NULL
   );
   ```

3. (Optional) Insert sample data for testing:
   ```sql
   -- Sample categories
   INSERT INTO categories (name, description) VALUES
   ('Beverages', 'Drinks and beverages'),
   ('Snacks', 'Snack items');

   -- Sample products
   INSERT INTO products (name, price, stock, category_id) VALUES
   ('Mineral Water', 5000, 100, 1),
   ('Coffee', 15000, 50, 1),
   ('Chips', 10000, 75, 2);
   ```

## API Documentation

This project uses OpenAPI (Swagger) 2.0 to document its endpoints and Scalar to provide a modern documentation UI.

### Viewing Documentation

Once the server is running, you can access the interactive API reference at:

- **Scalar UI**: [http://localhost:8080/reference](http://localhost:8080/reference) (Recommended - Modern, interactive UI)
- **Swagger JSON**: [http://localhost:8080/docs/swagger.json](http://localhost:8080/docs/swagger.json) (Raw specification)

The Scalar UI provides:
- Interactive API testing
- Request/response examples
- Schema documentation
- Authentication testing

### Generating Documentation

If you modify the API endpoints or add new godoc comments in handler files, you must regenerate the documentation:

1. **Install Swag CLI** (if not already installed):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate Docs**:
   Run the following command in the project root:
   ```bash
   swag init
   ```
   This will update the files in the `docs/` directory (`docs.go`, `swagger.json`, `swagger.yaml`).

3. **Restart the server** to see the updated documentation.

## Technology Stack

- **Language**: Go 1.25.6
- **Web Framework**: Native `net/http` (no external framework)
- **Database**: PostgreSQL
- **Database Driver**: pgx/v5 (PostgreSQL driver and toolkit)
- **Configuration**: Viper (environment and config file management)
- **API Documentation**: Swaggo (Swagger generator) with Scalar UI
- **Architecture**: Clean Architecture (Layered Pattern)

## Key Features Explained

### 1. Clean Architecture
The project follows a layered architecture pattern ensuring:
- **Separation of Concerns**: Each layer has a specific responsibility
- **Testability**: Business logic is isolated from infrastructure
- **Maintainability**: Changes in one layer don't affect others
- **Scalability**: Easy to add new features without breaking existing code

### 2. Transaction Management
The checkout system includes:
- Atomic transaction processing (all-or-nothing approach)
- Automatic stock reduction
- Bulk insert optimization for transaction details
- Complete audit trail of all transactions

### 3. Reporting System
Generate business insights with:
- Real-time sales summaries
- Best-selling product analysis
- Flexible date range filtering
- Revenue tracking

### 4. Database Connection Pooling
Optimized database performance with:
- Maximum 25 open connections
- Maximum 5 idle connections
- Automatic connection health checks

## Database Schema

The application uses the following main tables:

- **categories**: Stores product categories
  - `id`, `name`, `description`, `created_at`
- **products**: Stores product information
  - `id`, `name`, `price`, `stock`, `category_id`, `created_at`
- **transactions**: Stores transaction headers
  - `id`, `total_amount`, `created_at`
- **transaction_details**: Stores individual items in each transaction
  - `id`, `transaction_id`, `product_id`, `quantity`, `subtotal`

## Project Structure

```
kasir-api/
├── main.go                          # Application entry point, routing, and Swagger setup
├── go.mod                           # Go module definition
├── go.sum                           # Go module checksums
├── .env                             # Environment variables (not committed)
├── .gitignore                       # Git ignore rules
├── database/
│   └── database.go                  # Database connection and initialization
├── models/
│   ├── category.go                  # Category data model
│   ├── product.go                   # Product data model
│   ├── transaction.go               # Transaction and checkout models
│   └── report.go                    # Report data models
├── repositories/
│   ├── category_repository.go       # Category database operations
│   ├── product_repository.go        # Product database operations
│   ├── transaction_repository.go    # Transaction database operations
│   └── report_repository.go         # Report data queries
├── services/
│   ├── category_service.go          # Category business logic
│   ├── product_service.go           # Product business logic
│   ├── transaction_service.go       # Transaction business logic
│   └── report_service.go            # Report business logic
├── handlers/
│   ├── category_handler.go          # Category HTTP handlers
│   ├── product_handler.go           # Product HTTP handlers
│   ├── transaction_handler.go       # Transaction HTTP handlers
│   └── report_handler.go            # Report HTTP handlers
└── docs/                            # Auto-generated Swagger documentation
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

## Development Guide

### Adding a New Endpoint

To add a new feature, follow these steps:

1. **Define the Model** (`models/`):
   ```go
   type YourModel struct {
       ID        int       `json:"id"`
       Name      string    `json:"name"`
       CreatedAt time.Time `json:"created_at"`
   }
   ```

2. **Create the Repository** (`repositories/`):
   ```go
   type YourRepository struct {
       db *sql.DB
   }

   func (r *YourRepository) Create(model *YourModel) error {
       // Database operations
   }
   ```

3. **Create the Service** (`services/`):
   ```go
   type YourService struct {
       repo *repositories.YourRepository
   }

   func (s *YourService) CreateItem(model *YourModel) error {
       // Business logic
       return s.repo.Create(model)
   }
   ```

4. **Create the Handler** (`handlers/`):
   ```go
   type YourHandler struct {
       service *services.YourService
   }

   // Add godoc comments for Swagger
   // @Summary Create item
   // @Description Create a new item
   // @Tags items
   // @Accept json
   // @Produce json
   // @Param item body models.YourModel true "Item to create"
   // @Success 201 {object} models.YourModel
   // @Router /api/items [post]
   func (h *YourHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
       // Handle HTTP request/response
   }
   ```

5. **Register the Route** (`main.go`):
   ```go
   http.HandleFunc("/api/items", yourHandler.HandleItems)
   ```

6. **Regenerate Swagger Documentation**:
   ```bash
   swag init
   ```

### Running in Production

For production deployment, consider:

1. **Build the binary**:
   ```bash
   go build -o kasir-api main.go
   ```

2. **Use environment variables** for configuration (never commit `.env` file)

3. **Set up proper PostgreSQL connection** with SSL enabled

4. **Configure reverse proxy** (nginx/Apache) for the Go application

5. **Set up monitoring** and logging

## Recent Updates

### Latest Features (pertemuan-3 branch)

- ✅ **Sales Summary Report API**: Generate comprehensive sales reports with optional date filtering
  - Total revenue calculation
  - Transaction count tracking
  - Best-selling product identification

- ✅ **Transaction System Optimization**: Improved performance with bulk insert for transaction details

- ✅ **Enhanced Product Management**: Added name-based filtering for product search

### Previous Updates

- ✅ **Complete Transaction System**: Checkout functionality with automatic stock management
- ✅ **Database Integration**: Migrated from in-memory storage to PostgreSQL
- ✅ **Clean Architecture**: Implemented layered architecture pattern
- ✅ **API Documentation**: Added Swagger/OpenAPI documentation with Scalar UI
- ✅ **Configuration Management**: Environment-based configuration with Viper

## Environment Variables

| Variable  | Description                          | Example                                              |
|-----------|--------------------------------------|------------------------------------------------------|
| `PORT`    | Server port number                   | `8080`                                               |
| `DB_CONN` | PostgreSQL connection string         | `postgres://user:pass@localhost:5432/kasir_db?sslmode=disable` |

## Testing the API

You can test the API using various tools:

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Get all products
curl http://localhost:8080/api/products

# Create a new product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","price":25000,"stock":50,"category_id":1}'

# Checkout
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":2},{"product_id":2,"quantity":1}]}'

# Get sales report
curl "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-01-31"
```

### Using Scalar UI

Navigate to [http://localhost:8080/reference](http://localhost:8080/reference) for an interactive API testing interface.

## Troubleshooting

### Database Connection Issues

- Ensure PostgreSQL is running: `pg_isready`
- Check connection string format in `.env`
- Verify database exists: `psql -l`
- Check user permissions

### Port Already in Use

- Change the `PORT` in `.env` file
- Or kill the process using the port: `lsof -ti:8080 | xargs kill -9`

### Swagger Documentation Not Updating

- Regenerate docs: `swag init`
- Restart the server
- Clear browser cache

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is for educational purposes as part of "Jago Golang Dasar" learning program.

## Acknowledgments

- Built with Go's standard library and minimal dependencies
- PostgreSQL for reliable data persistence
- Swaggo for API documentation
- Scalar for modern API reference UI
