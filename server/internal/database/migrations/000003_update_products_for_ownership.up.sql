-- new columns
ALTER TABLE products
ADD COLUMN seller_id BIGINT NOT NULL;

ALTER TABLE products
ADD COLUMN price BIGINT NOT NULL;


-- new constraints
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
