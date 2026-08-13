ALTER TABLE products
DROP CONSTRAINT uq_products_seller_name;

ALTER TABLE products
DROP CONSTRAINT fk_products_seller;

ALTER TABLE products
DROP CONSTRAINT chk_products_price;

ALTER TABLE products
DROP COLUMN price;

ALTER TABLE products
DROP COLUMN seller_id;