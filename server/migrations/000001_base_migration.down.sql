-- order_items table
ALTER TABLE order_items
DROP CONSTRAINT chk_order_items_quantity;

ALTER TABLE order_items
DROP CONSTRAINT chk_order_items_price;

ALTER TABLE order_items
DROP CONSTRAINT uq_order_items_order_product;

ALTER TABLE order_items
DROP CONSTRAINT fk_order_items_product;

ALTER TABLE order_items
DROP CONSTRAINT fk_order_items_order;

DROP TABLE order_items;


-- orders table
ALTER TABLE orders
DROP CONSTRAINT chk_orders_amount;

ALTER TABLE orders
DROP CONSTRAINT chk_orders_payment;

ALTER TABLE orders
DROP CONSTRAINT fk_orders_user;

DROP TABLE orders;


-- cart_items table
ALTER TABLE cart_items
DROP CONSTRAINT chk_cart_items_quantity;

ALTER TABLE cart_items
DROP CONSTRAINT uq_cart_items_cart_product;

ALTER TABLE cart_items
DROP CONSTRAINT fk_cart_items_product;

ALTER TABLE cart_items
DROP CONSTRAINT fk_cart_items_cart;

DROP TABLE cart_items;


-- carts table
ALTER TABLE carts
DROP CONSTRAINT fk_carts_user;

DROP TABLE carts;


-- products table
ALTER TABLE products
DROP CONSTRAINT uq_products_seller_name;

ALTER TABLE products
DROP CONSTRAINT fk_products_seller;

ALTER TABLE products
DROP CONSTRAINT chk_products_price;

DROP TABLE products;


-- users table
DROP TABLE users;








