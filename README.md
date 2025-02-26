# Go E-Commerce API

This project provides a RESTful API for a basic e-commerce platform built using Go, Gin, and GORM.  It includes endpoints for user management, product catalog, categories, wishlists, and shopping carts.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
    - [User Management](#user-management)
    - [Product Management](#product-management)
    - [Category Management](#category-management)
    - [Wishlist Management](#wishlist-management)
    - [Cart Management](#cart-management)
- [Database](#database)
- [Error Handling](#error-handling)
- [Authentication/Authorization (TODO)](#authenticationauthorization-todo)
- [Input Validation](#input-validation)
- [Contributing](#contributing)
- [License](#license)

## Features

*   **User Management:** Create, list, get, update, and delete user accounts.
*   **Product Management:** Create, list, get, update, and delete products.
*   **Category Management:** Create, list, get, update, and delete product categories.
*   **Wishlist Management:** Add products to a user's wishlist, view a user's wishlist, view all wishlists and remove items from a wishlist.
*   **Cart Management:** Add products to a user's shopping cart, view a user's cart, update quantities of items in the cart, and remove items from the cart.
*   **Database Seeding:** Seed initial categories and products for development/testing.
*   **Automatic Database Migrations:** GORM automatically migrates the database schema based on the defined models.

## Technologies Used

*   **Go:** Programming language
*   **Gin:** Web framework
*   **GORM:** ORM (Object-Relational Mapper) for database interaction
*   **MySQL:** Database (but easily adaptable to other databases supported by GORM)
*   **godotenv:** For loading environment variables from a `.env` file

## Prerequisites

*   Go installed (version 1.18 or higher)
*   MySQL database server running
*   `make` command installed (optional, but recommended for convenience)

## Setup

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd <project_directory>
    ```

2.  **Install dependencies:**

    ```bash
    go mod download
    ```

3.  **Create a `.env` file:**

    Copy the example `.env.example` file to `.env` and configure the database connection settings:

    ```bash
    cp .env.example .env
    ```

    Edit `.env` and fill in the appropriate values.

## Environment Variables

The following environment variables are required:

*   `PORT`: The port the application will listen on (e.g., `8080`).
*   `DB_USER`: The MySQL database user.
*   `DB_PASSWORD`: The MySQL database password.
*   `DB_HOST`: The MySQL database host (e.g., `localhost`).
*   `DB_PORT`: The MySQL database port (e.g., `3306`).
*   `DB_NAME`: The MySQL database name.
*   `SEED_PRODUCT_COUNT`: The number of products to seed when the application starts (optional).

## Running the Application

1.  **Initialize the database:**

    The application automatically runs database migrations on startup.  Ensure your MySQL server is running and the connection details in `.env` are correct.

2.  **Run the application:**

    ```bash
    go run main.go
    ```

    The application will start, and you can access the API endpoints.

## API Endpoints

All API endpoints return JSON responses.

### User Management

*   **`POST /api/users`:** Create a new user.  Requires JSON body with user details (firstName, lastName, email, password, phone, address).
*   **`GET /api/users`:** List all users.
*   **`GET /api/users/:userid`:** Get a single user by ID.
*   **`PATCH /api/users/:userid`:** Update an existing user.  Requires JSON body with the fields to update.
*   **`DELETE /api/users/:userid`:** Delete a user.

### Product Management

*   **`POST /api/product`:** Create a new product.  Requires JSON body with product details (SKU, name, description, price, categoryID).
*   **`GET /api/product`:** List all products.
*   **`GET /api/product/:productid/`:** Get a single product by ID.
*   **`PATCH /api/product/:productid/`:** Update an existing product. Requires JSON body with the fields to update.
*   **`DELETE /api/product/:productid/`:** Delete a product.

### Category Management

*   **`POST /api/categories`:** Create a new category.  Requires JSON body with category details (name, description).
*   **`GET /api/categories`:** List all categories.
*   **`GET /api/categories/:categoryid/`:** Get a single category by ID.
*   **`PATCH /api/categories/:categoryid/`:** Update an existing category.  Requires JSON body with the fields to update (name, description).
*   **`DELETE /api/categories/:categoryid/`:** Delete a category.

### Wishlist Management

*   **`POST /api/wishlists`:** Add a product to a user's wishlist.  Requires JSON body with `userID` and `productID`.
*   **`GET /api/wishlists/:userid/`:** View a user's wishlist.
*   **`GET /api/wishlists`:** View all wishlists.
*   **`DELETE /api/wishlists/:wishlistid/`:** Remove an item from a user's wishlist.

### Cart Management

*   **`POST /api/carts`:** Add a product to a user's shopping cart. Requires JSON body with `userID`, `productID`, and `quantity`.
*   **`GET /api/carts/:userid/`:** View a user's shopping cart.
*   **`PATCH /api/carts/:cartid/`:** Update the quantity of an item in a user's shopping cart. Requires JSON body with `quantity`.
*   **`DELETE /api/carts/:cartid/`:** Remove an item from a user's shopping cart.

## Database

The application uses a MySQL database.  The database schema is automatically created and updated by GORM based on the model definitions in `models/models.go`.  The following tables are created:

*   `users`
*   `products`
*   `categories`
*   `wishlists`
*   `carts`

## Error Handling

The API uses a consistent error handling pattern:

*   **`400 Bad Request`:** For invalid requests (e.g., missing required fields, invalid data types).
*   **`404 Not Found`:** For requests to non-existent resources.
*   **`500 Internal Server Error`:** For unexpected server errors.

Error messages are returned in JSON format with a `message` field.

## Authentication/Authorization (TODO)

**Important:** Currently, there is NO authentication or authorization implemented in this project. All API endpoints are publicly accessible.

Future development should include:

*   **User Authentication:** Implement user login/registration using a secure authentication mechanism (e.g., JWTs - JSON Web Tokens).
*   **Authorization:** Implement authorization rules to restrict access to resources based on user roles and permissions.

## Input Validation

The API uses Gin's binding and validation features to validate input data.  The `binding:"required"` tag is used to mark required fields.  Custom validation logic can be added as needed.

## License

This project is licensed under the [MIT License](LICENSE).
