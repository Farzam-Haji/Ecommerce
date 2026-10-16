# E-commerce API



## Overview
>This project is a e-commerce backend built with Go, Gin, and PostgreSQL. It provides user authentication, product management, shopping carts, order management, and payment integration. The API uses JWT authentication for protected endpoints and PostgreSQL migrations for database setup.

### Base URL
```text
/api/v1
```
### Units
> Prices and order amounts are represented in minor currency units.
`Example:` `1000 = 10.00`



## Authentication
> Protected endpoints require a JWT.
```text
Authorization: Bearer <jwt>
```
> JWTs expire after 72 hours.
#### Errors
- **401 Unauthorized** - `ErrUnauthorized`: `authentication required`
- **401 Unauthorized** - `ErrInvalidToken`: `invalid token`
- **401 Unauthorized** - `ErrExpiredToken`: `token has expired`



## Error Responses
All error responses use:
```json
{
    "error": "error message"
}
```




## Configuration

>### Environment

The application reads its configuration from environment variables.

Create a `.env` file in the project root:

```env
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DBNAME=your_database_name

JWT_SECRET=your_jwt_secret

SERVER=localhost:8080
BANK_URL=http://localhost:8081/api/v1/payments
```
| Variable | Description |
| -------- | ----------- |
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_HOST` | PostgreSQL host |
| `POSTGRES_PORT` | PostgreSQL port |
| `POSTGRES_DBNAME` | PostgreSQL database name |
| `JWT_SECRET` | Secret used to sign JWTs |
| `SERVER` | Address where the API server listens |
| `BANK_URL` | Base URL of the payment service |


>### Dependencies

This project uses Go modules and dependencies are defined in `go.mod`.

To download them:
```bash
go mod download
```

>### Database Setup

This project uses PostgreSQL and `golang-migrate` for database migrations.

#### Requirements

- Go
- PostgreSQL
- golang-migrate

#### Create the Database

Create a PostgreSQL database matching `POSTGRES_DBNAME` in `.env`.

#### Run Migrations

From the project root:

```bash
migrate -path ./migrations \
  -database "postgres://<USER>:<PASSWORD>@<HOST>:<PORT>/<DBNAME>?sslmode=disable" \
  up
```

And to rollback the latest migration :

```bash
migrate -path ./migrations \
  -database "postgres://<USER>:<PASSWORD>@<HOST>:<PORT>/<DBNAME>?sslmode=disable" \
  down 1
```

>### Start Server

The API server starts at the address configured by `SERVER` in `.env`.
To run the server from the project root:
```bash
go run ./cmd/api
```

>### Bank Server

Order creation and payment functionality depend on a separate bank/payment service.
For development and testing, you can use the [mock bank](https://github.com/Farzam-Haji/mock-bank) I created for this project.
Clone and run the mock bank server before testing order payment functionality. Make sure its address matches the `BANK_URL` configured in `.env`.









--------------------------------
## Authentication Endpoints



### POST /register
- Create a new user.

>#### Authentication
- Not required.

>#### Request
```json
{
    "email": "user@example.com",
    "password": "example-password"
}
```
- `email`: required string, must be a valid email address, must be unique
- `password`: required string, minimum 8 characters

>#### Response
- **201 Created** - `json body:`
```json
{
    "user": "user@example.com",
    "token": "<jwt>"
}
```

>#### Errors
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrEmailUnavailable`: `email is already registered`
- **400 Bad Request** - `ErrInvalidPassword`: `password must contain at least 8 characters`
- **500 Internal Server Error** - `ErrTokenGenFailed`: `failed to generate authentication token`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### POST /login
- Generate JWT token for a user if given a valid email and password.

>#### Authentication
- Not required.

>#### Request
```json
{
    "email": "user@example.com",
    "password": "example-password"
}
```
- `email`: required string, must be a valid email address
- `password`: required string, minimum 8 characters

>#### Response
- **200 OK** - `json body:`
```json
{
    "token": "<JWT>"
}
```

>#### Errors
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **401 Unauthorized** - `ErrWrongCredential`: `one or both credentials are wrong`
- **500 Internal Server Error** - `ErrTokenGenFailed`: `failed to generate authentication token`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### GET /me
- Returns the current authentication state and, if authenticated, the user's email.

>#### Authentication
- Not required.

>#### Request
- No Request body.

>#### Responses
- **200 OK** - `json body:`
```json
{
    "authenticated": true,
    "user": "user@example.com"
}
```
- `authenticated`: Boolean indicating whether the request has a valid JWT.
- `user`: Present only when authenticated is true; contains the authenticated user's email.

- When the request is unauthenticated, the response is:
```json
{
    "authenticated": false
}
```




-------------------------------
## Product Endpoints



### GET /products
- Returns all registered products.

>#### Authentication
- Not required.

>#### Request
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "products": [
        {
          "id": 1,
          "name": "apple",
          "quantity": 10,
          "seller_id": 1,
          "price": 50
        },
        {
          "id": 2,
          "name": "orange",
          "quantity": 100,
          "seller_id": 2,
          "price": 55
        },
        {
          "id": 3,
          "name": "Peach",
          "quantity": 10,
          "seller_id": 1,
          "price": 100
        }
    ]
}
```

>#### Errors
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### GET /products/me
- Returns all products registered by the authenticated user.

>#### Authentication
- Required.

>#### Request
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "products": [
        {
          "id": 1,
          "name": "apple",
          "quantity": 10,
          "seller_id": 1,
          "price": 50
        },
        {
          "id": 2,
          "name": "orange",
          "quantity": 100,
          "seller_id": 1,
          "price": 55
        }
    ]
}
```

>#### Errors
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### POST /products
- Creates a new product owned by the authenticated user.

>#### Authentication
- Required.

>#### Request
```json
{
    "name": "apple",
    "quantity": 1000,
    "price": 50
}
```
- `name`: required string
- `quantity`: required integer, value must be greater than or equal to 0
- `price`: required integer, value must be greater than or equal to 0

#### Response
- **201 Created** - `json body:`
```json
{
    "newProduct": {
      "id": 1,
      "name": "apple",
      "quantity": 10,
      "seller_id": 1,
      "price": 50
    }
}
```
>#### Errors
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrInvalidName`: `invalid product name`
- **400 Bad Request** - `ErrInvalidQuantity`: `invalid quantity`
- **400 Bad Request** - `ErrInvalidPrice`: `invalid price`
- **409 Conflict** - `ErrDuplicateProduct`: `product already exists`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### GET /products/:id
- Fetch a product by its ID.

>#### Authentication
- Not required

>#### Request
- **Path parameter** - `id`: required integer, Product ID
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "product": {
      "id": 1,
      "name": "apple",
      "quantity": 10,
      "seller_id": 1,
      "price": 50
    }
}
```

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **404 Not Found** - `ErrProductNotFound`: `product not found`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### PUT /products/:id
- Updates all fields of a product owned by the authenticated user.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: required integer, Product ID
```json
{
    "name": "apple",
    "quantity": 10,
    "price": 100
}
```
- `name`: required string
- `quantity`: required integer, must be greater than or equal to 0
- `price`: required integer, must be greater than or equal to 0

>#### Response
- **200 OK**

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **400 Bad Request** - `ErrInvalidName`: `invalid product name`
- **400 Bad Request** - `ErrInvalidQuantity`: `invalid quantity`
- **400 Bad Request** - `ErrInvalidPrice`: `invalid price`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### PATCH /products/:id
- Updates the provided fields of a product owned by the authenticated user.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: required integer, Product ID
```json
{
    "name": "apple",
    "quantity": 10,
    "price": 10
}
```
- `name`: optional string
- `quantity`: optional integer, must be greater than or equal to 0
- `price`: optional integer, must be greater than or equal to 0
- All fields are optional. Only provided fields are updated.

>#### Response
- **200 OK**

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **400 Bad Request** - `ErrInvalidName`: `invalid product name`
- **400 Bad Request** - `ErrInvalidQuantity`: `invalid quantity`
- **400 Bad Request** - `ErrInvalidPrice`: `invalid price`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### DELETE /products/:id
- Delete a product owned by authenticated user.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: required integer, Product ID
- No request body.

>#### Response
- **204 No Content**

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`






---------------------------------------------
## Cart Endpoints
### GET /cart
- Returns the cart of the authenticated user.
- Carts are created on demand when an authenticated user adds their first item to the cart.

>#### Authentication
- Required.

>#### Request
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "items": [
      {
        "id": 1,
        "cart_id": 1,
        "product_id": 1,
        "name": "apple",
        "quantity": 20,
        "price": 1000
      },
      {
        "id": 2,
        "cart_id": 1,
        "product_id": 2,
        "name": "orange",
        "quantity": 100,
        "price": 5500
      }
    ],
    "notification": "some items had been removed from your cart",
    "deleted_items": ["peach", "pear"]
}
```
- `items`: array of cart items. Returns `[]` when the cart is empty.
- `notification`: included only when one or more products were removed from the cart because they were no longer available.
-`deleted_items`: array containing the names of products that were removed. Included only when one or more products were removed.
- `item.price`: total price of the cart item, calculated as the product's unit price multiplied by the item's quantity.

>#### Errors
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### DELETE /cart
- Delete the cart of authenticated user.

>#### Authentication
- Required.

>#### Request
- No request body.

>#### Response
- **204 No Content** 

>#### Errors
- **404 Not Found** - `ErrCartNotFound`: `cart not found`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### POST /cart/items
- Adds an item to the authenticated user's cart. 
- If the product is already in the cart, its quantity is increased.

>#### Authentication
- Required.

>#### Request
```json
{
    "product_id": 123,
    "quantity": 80
}
```
- `product_id`: required integer, must be greater than 0.
- `quantity`: optional integer, must be greater than 0. Defaults to `1` when omitted.

>#### Response
- **201 Created** - `json body:`
```json
{
    "item": {
      "id": 1,
      "cart_id": 1,
      "product_id": 123,
      "name": "apple",
      "quantity": 80,
      "price": 8000
    }
}
```
- Cart item `price` is total price for the item.

>#### Errors
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **400 Bad Request** - `ErrInvalidQuantity`: `invalid quantity`
- **404 Not Found** - `ErrProductNotFound`: `product not found`
- **409 Conflict** - `ErrInsufficientQuantity`: `product is not available in the requested quantity`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`


### PATCH /cart/items/:id
- Updates the quantity of an item in the authenticated user's cart.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: required integer, Cart item ID
```json
{
  "quantity": 12
}
```
- `quantity`: required integer, must be greater than or equal to 0.
- If `quantity` is `0` item will be removed instead.

>#### Response
- `If the item remains in the cart:`**200 OK** - `json body:`
```json
{
    "item": {
        "id":1,
        "cart_id": 1,
        "product_id": 1,
        "name": "apple",
        "quantity": 12,
        "price": 1200
    }
}
```
- `If quantity is 0 and the item is removed:`**204 No Content** 

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **400 Bad Request** - `ErrInvalidQuantity`: `invalid quantity`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **404 Not Found** - `ErrItemNotFound`: `cart item not found`
- **404 Not Found** - `ErrProductNotFound`: `product not found`
- **409 Conflict** - `ErrInsufficientQuantity`: `product is not available in the requested quantity`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### DELETE /cart/items/:id
- Deletes an item from the authenticated user's cart.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: Required integer. Cart item ID
- No request body.

>#### Response
- **204 No Content**

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **400 Bad Request** - `ErrInvalidID`: `invalid id`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **404 Not Found** - `ErrItemNotFound`: `cart item not found`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`







-------------------------------
## Order Endpoints



### POST /orders
- Create a new order based on authenticated user's cart.
- Pending orders expire after 10 minutes. While an order is pending, its items' quantities are reserved.
- Each user can have at most one pending order.
- Repeating the request while a pending order exists refreshes its expiration time. Changes made to the cart after the order is created do not affect the pending order.

>#### Authentication
- Required.

>#### Request
- No request body.

>### Response
- **200 OK** - `json body:`
```json
{
    "order": {
        "id":1,
            "user_id":1,
            "payment":"pending",
            "payment_id": null,
            "amount":5050,
            "created_at":"2026-09-15T02:05:35.927359+03:30",
            "expires_at":"2026-09-15T02:15:35.929027+03:30",
            "updated_at":"2026-09-15T02:05:35.927359+03:30"

    },
    "payment_url":"http://localhost:8081/payment/123"
}
```
- `payment_url` is where payment may happen.

>#### Errors
- **400 Bad Request** - `ErrEmptyOrder`: `cart is empty`
- **404 Not Found** - `ErrCartNotFound`: `cart not found`
- **404 Not Found** - `ErrProductNotFound`: `product not found`
- **409 Conflict** - `ErrInsufficientQuantity`: `product is not available in the requested quantity`
- **500 Internal Server Error** - `ErrRequestBank`: `failed to create payment session`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### GET /orders
- Returns the order history of authenticated user.

>#### Authentication
- Required.

>#### Request
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "orders": [
        {
            "id":1,
            "user_id":1,
            "payment":"pending",
            "payment_id": null,
            "amount":5050,
            "created_at":"2026-09-15T02:05:35.927359+03:30",
            "expires_at":"2026-09-15T02:15:35.929027+03:30",
            "updated_at":"2026-09-15T02:05:35.927359+03:30"
        },
        {
            "id":2,
            "user_id":1,
            "payment":"paid",
            "payment_id": 123,
            "amount":6300,
            "created_at":"2026-09-15T03:05:35.927359+03:30",
            "expires_at":"2026-09-15T03:15:35.929027+03:30",
            "updated_at":"2026-09-15T03:10:35.463994+03:30"
        }
    ]
}
```
- `payment`: Payment status. Possible values are `pending`, `paid`, `failed`, `expired` and `canceled`.
- `payment_id`: The payment ID returned by the bank when the payment result is received. null until the bank sends the payment result.

>#### Errors
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### GET /orders/:id
- Returns an order belonging to the authenticated user and its items.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: Required integer. Order ID
- No request body.

>#### Response
- **200 OK** - `json body:`
```json
{
    "order": {
        "id":3,
        "user_id":2,
        "payment":"expired",
        "payment_id": null,
        "amount":87400,
        "created_at":"2026-09-15T02:05:35.927359+03:30",
        "expires_at":"2026-09-15T02:15:35.929027+03:30",
        "updated_at":"2026-09-15T02:15:40.927359+03:30"
    },
    "items": [
        {
            "id":21,
            "order_id":3,
            "product_id":44,
            "name":"pencil",
            "price":800,
            "quantity":100
        },
        {
            "id":22,
            "order_id":3,
            "product_id":786,
            "name":"pen",
            "price":3700,
            "quantity":2
        }
    ]
}
```
- `product_id` is `null` if the referenced product has been deleted. The product name and price stored in the order item remain available in the order history.

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **404 Not Found** - `ErrOrderNotFound`: `order not found`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



### POST /orders/:id/cancel 
- Cancels the authenticated user's pending order.
- The order is retained in order history with a canceled payment status.

>#### Authentication
- Required.

>#### Request
- **Path parameter** - `id`: Required integer. Order ID
- No request body.

>#### Response
- **200 OK**

>#### Errors
- **400 Bad Request** - `ErrInvalidParam`: `invalid path parameter`
- **403 Forbidden** - `ErrForbidden`: `authenticated user is not allowed to perform this action`
- **404 Not Found** - `ErrOrderNotFound`: `order not found`
- **404 Not Found** - `ErrNonPending`: `order is not pending`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`




------------------------------------------
## Payment Callback
### POST /orders/result
- Receives the payment result from the bank and updates the corresponding order.
- This endpoint is intended for server-to-server communication from the bank.

>#### Authentication
- Not required.

>#### Request
```json
{
    "payment_id":123,
    "order_id":2,
    "status":"paid"
}
```
- `payment_id`: Payment ID assigned by the bank.
- `order_id`: ID of the order associated with the payment.
- `status`: Payment result. Possible values are `paid` and `failed`.

>#### Response
- **200 OK**
    - Payment result was received and the order was updated.
    - The endpoint does not redirect to the order endpoint. The client can retrieve the updated order through `GET /orders/:id`.

>#### Errors
- **400 Bad Request** - `ErrInvalidJSON`: `invalid request body`
- **500 Internal Server Error** - `ErrInternal`: `internal server error`



--------------------------------------------------------------------------------
>## Order lifecycle
After creating an order, the client should redirect the user to `payment_url`. 
If the user chooses to leave the payment flow without completing payment, the client may cancel the pending order using `POST /orders/:id/cancel`.
Canceling an order releases its reserved product quantities. The canceled order remains in the user's order history with a canceled payment status.
This flow keeps order and cart synced, and avoid unsync expiration period between bank and api.
