ALTER TABLE cart_items
DROP CONSTRAINT chk_cart_items_quantity;

ALTER TABLE cart_items
DROP CONSTRAINT uq_cart_items_cart_product;

ALTER TABLE cart_items
DROP CONSTRAINT fk_cart_items_product;

ALTER TABLE cart_items
DROP CONSTRAINT fk_cart_items_cart;

DROP TABLE cart_items;


ALTER TABLE carts
DROP CONSTRAINT fk_carts_user;

DROP TABLE carts;