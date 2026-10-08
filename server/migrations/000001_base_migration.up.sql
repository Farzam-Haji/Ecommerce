-- user
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);




-- product
CREATE TABLE products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    seller_id BIGINT NOT NULL,
    price BIGINT NOT NULL
);

ALTER TABLE products 
ADD CONSTRAINT chk_products_price 
CHECK (price >= 0);

ALTER TABLE products
ADD CONSTRAINT fk_products_seller
FOREIGN KEY (seller_id) REFERENCES users(id)
ON DELETE CASCADE;

ALTER TABLE products
ADD CONSTRAINT uq_products_seller_name
UNIQUE (seller_id, name);




-- cart
CREATE TABLE carts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE carts
ADD CONSTRAINT fk_carts_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;



CREATE TABLE cart_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cart_id BIGINT NOT NULL,
    product_id BIGINT,
    name TEXT NOT NULL,
    quantity INT NOT NULL
);

ALTER TABLE cart_items
ADD CONSTRAINT fk_cart_items_cart
FOREIGN KEY (cart_id) REFERENCES carts(id)
ON DELETE CASCADE;

ALTER TABLE cart_items
ADD CONSTRAINT fk_cart_items_product
FOREIGN KEY (product_id) REFERENCES products(id)
ON DELETE SET NULL;

ALTER TABLE cart_items
ADD CONSTRAINT uq_cart_items_cart_product
UNIQUE (cart_id, product_id);

ALTER TABLE cart_items
ADD CONSTRAINT chk_cart_items_quantity
CHECK (quantity > 0);




-- order
CREATE TABLE orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT,
    payment TEXT NOT NULL DEFAULT 'pending',
    payment_id BIGINT UNIQUE,
    amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

);

ALTER TABLE orders
ADD CONSTRAINT fk_orders_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE SET NULL;

ALTER TABLE orders
ADD CONSTRAINT chk_orders_payment
CHECK (payment IN ('pending', 'paid', 'failed', 'expired', 'canceled'));

ALTER TABLE orders
ADD CONSTRAINT chk_orders_amount
CHECK (amount >= 0);



CREATE TABLE order_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT,
    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    quantity INT NOT NULL
);

ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_order
FOREIGN KEY (order_id) REFERENCES orders(id)
ON DELETE CASCADE;

ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_product
FOREIGN KEY (product_id) REFERENCES products(id)
ON DELETE SET NULL;

ALTER TABLE order_items
ADD CONSTRAINT uq_order_items_order_product
UNIQUE (order_id, product_id);

ALTER TABLE order_items
ADD CONSTRAINT chk_order_items_price
CHECK (price >= 0);

ALTER TABLE order_items
ADD CONSTRAINT chk_order_items_quantity
CHECK (quantity > 0);


