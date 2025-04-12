CREATE TABLE customers (
                           id INTEGER PRIMARY KEY AUTOINCREMENT,
                           created_at DATETIME DEFAULT (datetime('now')),
                           updated_at DATETIME DEFAULT (datetime('now')),
                           deleted_at DATETIME DEFAULT NULL,
                           first_name TEXT NOT NULL,
                           last_name TEXT NOT NULL,
                           address TEXT NOT NULL,
                           email TEXT NOT NULL UNIQUE
);

CREATE TABLE orders (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        created_at DATETIME DEFAULT (datetime('now')),
                        updated_at DATETIME DEFAULT (datetime('now')),
                        deleted_at DATETIME DEFAULT NULL,
                        type TEXT NOT NULL,
                        size TEXT NOT NULL,
                        quantity INTEGER NOT NULL,
                        customer_id INTEGER NOT NULL,
                        status TEXT NOT NULL DEFAULT 'new'
);

CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
