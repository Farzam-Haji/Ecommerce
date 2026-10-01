ALTER TABLE orders
ADD COLUMN payment_id BIGINT UNIQUE;

ALTER TABLE orders
ADD COLUMN amount BIGINT NOT NULL DEFAULT 0;

ALTER TABLE orders
ADD CONSTRAINT chk_orders_amount
CHECK (amount >= 0);

ALTER TABLE orders
DROP CONSTRAINT chk_orders_payment;

ALTER TABLE orders
ADD CONSTRAINT chk_orders_payment
CHECK (payment IN ('pending', 'paid', 'failed', 'expired', 'canceled'));