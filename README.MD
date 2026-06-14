# Rentalin API

A simple rental management REST API built with Go.

This project was created for learning and portfolio purposes, focusing on backend development concepts such as authentication, authorization, database transactions, payment gateway integration, and clean application architecture.

Feel free to explore, use, and modify this project as part of your learning journey. Hope it helps and inspires you to build something awesome!

## Features

* JWT Authentication
* Role-Based Access Control
* User Management
* Product Management
* Rental Management
* Xendit Payment Integration
* Dashboard Statistics
* Structured Logging with Logrus

## Tech Stack

* Go (Golang)
* Echo
* PostgreSQL
* PGX
* JWT
* Logrus
* Xendit

## Database Schema (ERD)

```mermaid
erDiagram
    USERS {
        bigint id PK
        varchar username
        varchar email
        varchar phone
        varchar address
        text password
        user_role role
        timestamp created_at
        timestamp updated_at
    }

    PRODUCTS {
        bigint id PK
        varchar name
        bigint daily_price
        int stock
        timestamp created_at
        timestamp updated_at
    }

    RENTALS {
        bigint id PK
        bigint customer_id FK
        bigint product_id FK
        bigint created_by FK
        date start_date
        date end_date
        bigint total_price
        rental_status status
        timestamp created_at
        timestamp updated_at
    }

    PAYMENTS {
        bigint id PK
        bigint customer_id FK
        bigint rental_id FK
        varchar external_id
        text invoice_url
        bigint amount
        bigint paid_amount
        varchar currency
        varchar method
        varchar payment_channel
        payment_status status
        timestamp expired_at
        timestamp paid_at
        text description
        jsonb callback_payload
        timestamp created_at
        timestamp updated_at
    }

    USERS ||--o{ RENTALS : customer
    USERS ||--o{ RENTALS : created_by
    PRODUCTS ||--o{ RENTALS : rented_in

    USERS ||--o{ PAYMENTS : pays
    RENTALS ||--|| PAYMENTS : payment
```

## Getting Started

Clone the repository:

```bash
git clone https://github.com/yourusername/rentalin.git
cd rentalin
```

Copy the example environment file:

```bash
cp .env.example .env
```

Then update the values in `.env` according to your local environment.

> All required environment variables are already documented in `.env.example`.

Run the application:

```bash
go run cmd/server/main.go
```

The server will start on:

```text
http://localhost:8080
```

## API Collection

Postman collection is available in:

```text
docs/rentalin-api.postman_collection.json
```
