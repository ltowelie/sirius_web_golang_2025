DROP INDEX IF EXISTS idx_orders_customer_id;
DROP INDEX IF EXISTS idx_orders_deleted_at;
DROP INDEX IF EXISTS idx_customers_deleted_at;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customers;