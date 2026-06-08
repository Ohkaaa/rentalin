-- ENUM TYPES
CREATE TYPE user_role AS ENUM (
    'admin',
    'customer'
);

CREATE TYPE rental_status AS ENUM (
    'pending',
    'ongoing',
    'completed',
    'overdue',
    'cancelled'
);

CREATE TYPE payment_status AS ENUM (
    'pending',
    'paid',
    'failed',
    'expired'
);

-- USERS
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(150) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- PRODUCTS
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    daily_price BIGINT NOT NULL CHECK (daily_price >= 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- RENTALS
CREATE TABLE rentals (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL, -- admin ID or user ID
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_price BIGINT NOT NULL CHECK (total_price >= 0),
    status rental_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_rentals_customer FOREIGN KEY (customer_id) REFERENCES users(id),
    CONSTRAINT fk_rentals_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_rentals_admin FOREIGN KEY (created_by) REFERENCES users(id),

    CHECK (end_date >= start_date)
);

CREATE INDEX idx_rentals_customer ON rentals(customer_id);
CREATE INDEX idx_rentals_product ON rentals(product_id);
CREATE INDEX idx_rentals_status ON rentals(status);

-- PAYMENTS
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    rental_id BIGINT NOT NULL,
    external_id VARCHAR(100) NOT NULL UNIQUE,
    invoice_url TEXT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount >= 0),
    paid_amount BIGINT,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    method VARCHAR(50),           -- e.g: BANK_TRANSFER
    payment_channel VARCHAR(50),  -- e.g: BCA, OVO
    status payment_status NOT NULL DEFAULT 'pending',
    expired_at TIMESTAMP NOT NULL,
    paid_at TIMESTAMP NULL,
    description TEXT,
    callback_payload JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_payments_customer FOREIGN KEY (customer_id)
        REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_payments_rental FOREIGN KEY (rental_id)
        REFERENCES rentals(id) ON DELETE CASCADE,

    CONSTRAINT unique_rental_payment UNIQUE (rental_id)
);

CREATE INDEX idx_payments_customer ON payments(customer_id);
CREATE INDEX idx_payments_rental ON payments(rental_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_external_id ON payments(external_id);