CREATE TABLE orders (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id BIGINT,
	payment TEXT NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	expires_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE orders
ADD CONSTRAINT fk_orders_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE SET NULL;

ALTER TABLE orders
ADD CONSTRAINT chk_orders_payment
CHECK (payment IN ('pending', 'paid', 'failed', 'expired'));


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
FOREIGN KEY	(order_id) REFERENCES orders(id)
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