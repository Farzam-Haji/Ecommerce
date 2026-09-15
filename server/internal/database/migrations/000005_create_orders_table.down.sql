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


ALTER TABLE orders
DROP CONSTRAINT chk_orders_payment;

ALTER TABLE orders
DROP CONSTRAINT fk_orders_user;

DROP TABLE orders;