CREATE TABLE carts (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id BIGINT NOT NULL UNIQUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

ALTER TABLE	cart_items
ADD CONSTRAINT chk_cart_items_quantity
CHECK (quantity > 0);