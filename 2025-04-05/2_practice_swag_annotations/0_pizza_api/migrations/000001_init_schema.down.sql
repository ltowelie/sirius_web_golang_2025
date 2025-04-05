DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_orders_customer_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customers;