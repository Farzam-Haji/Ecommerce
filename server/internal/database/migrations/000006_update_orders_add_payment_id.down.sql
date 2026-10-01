ALTER TABLE orders
DROP COLUMN payment_id;

ALTER TABLE orders
DROP CONSTRAINT chk_orders_amount;

ALTER TABLE orders
DROP COLUMN amount;