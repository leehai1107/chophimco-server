/* =========================================================
DATABASE: Mechanical Keyboard Shop (FULL)
Engine   : PostgreSQL
========================================================= */

-- =======================
-- 1. CREATE DATABASE
-- =======================
CREATE DATABASE mechanical_keyboard_shop;

\c mechanical_keyboard_shop;

-- =======================
-- 2. ROLES
-- =======================
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO roles (name) VALUES ('admin'), ('customer');

-- =======================
-- 3. USERS
-- =======================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL REFERENCES roles (id),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 4. CATEGORIES
-- =======================
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- =======================
-- 5. BRANDS
-- =======================
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- =======================
-- 6. SWITCHES
-- =======================
CREATE TABLE switches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(50), -- Linear, Tactile, Clicky
    brand VARCHAR(100)
);

-- =======================
-- 7. PRODUCTS
-- =======================
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    category_id INT REFERENCES categories (id),
    brand_id INT REFERENCES brands (id),
    description TEXT,
    base_price DECIMAL(12, 2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 8. PRODUCT VARIANTS
-- =======================
CREATE TABLE product_variants (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    switch_id INT REFERENCES switches (id),
    layout VARCHAR(20), -- 60%, 65%, TKL, Fullsize
    connection_type VARCHAR(50), -- Wired, Wireless, Bluetooth
    hotswap BOOLEAN DEFAULT FALSE,
    led_type VARCHAR(50), -- RGB, White
    price DECIMAL(12, 2) NOT NULL,
    stock INT DEFAULT 0,
    sku VARCHAR(100) UNIQUE
);

-- =======================
-- 9. CARTS
-- =======================
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 10. CART ITEMS
-- =======================
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
    product_variant_id INT NOT NULL REFERENCES product_variants (id),
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_variant_id)
);

-- =======================
-- 11. VOUCHERS
-- =======================
CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) NOT NULL, -- percent | fixed
    discount_value DECIMAL(12, 2) NOT NULL,
    min_order_value DECIMAL(12, 2) DEFAULT 0,
    max_discount_value DECIMAL(12, 2),
    usage_limit INT,
    usage_per_user INT DEFAULT 1,
    used_count INT DEFAULT 0,
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 12. USER VOUCHERS
-- =======================
CREATE TABLE user_vouchers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    voucher_id INT NOT NULL REFERENCES vouchers (id) ON DELETE CASCADE,
    used_count INT DEFAULT 0,
    UNIQUE (user_id, voucher_id)
);

-- =======================
-- 13. PRODUCT DISCOUNTS
-- =======================
CREATE TABLE product_discounts (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    discount_type VARCHAR(20) NOT NULL, -- percent | fixed
    discount_value DECIMAL(12, 2) NOT NULL,
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 14. ORDERS
-- =======================
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id),
    voucher_id INT REFERENCES vouchers (id),
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    total_amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending, paid, shipped, completed, cancelled
    shipping_address TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =======================
-- 15. ORDER ITEMS
-- =======================
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_variant_id INT NOT NULL REFERENCES product_variants (id),
    price DECIMAL(12, 2) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0)
);

-- =======================
-- 16. PAYMENTS
-- =======================
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    payment_method VARCHAR(50), -- COD, Momo, VNPay
    payment_status VARCHAR(50), -- pending, success, failed
    paid_at TIMESTAMP
);

-- =======================
-- 17. REVIEWS
-- =======================
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id),
    product_id INT NOT NULL REFERENCES products (id),
    rating INT CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);

-- =======================
-- 18. INDEXES (PERFORMANCE)
-- =======================
CREATE INDEX idx_products_category ON products (category_id);

CREATE INDEX idx_products_brand ON products (brand_id);

CREATE INDEX idx_variants_product ON product_variants (product_id);

CREATE INDEX idx_variants_switch ON product_variants (switch_id);

CREATE INDEX idx_orders_user ON orders (user_id);

CREATE INDEX idx_orders_voucher ON orders (voucher_id);

CREATE INDEX idx_voucher_code ON vouchers (code);

CREATE INDEX idx_voucher_active ON vouchers (is_active);

CREATE INDEX idx_user_voucher ON user_vouchers (user_id, voucher_id);

CREATE INDEX idx_product_discount ON product_discounts (product_id);

CREATE INDEX idx_reviews_product ON reviews (product_id);

-- =======================
-- END OF FILE
-- =======================