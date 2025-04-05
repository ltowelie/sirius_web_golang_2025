CREATE TABLE customers (
                           id INTEGER PRIMARY KEY AUTOINCREMENT,
                           created_at DATETIME DEFAULT (datetime('now')),
                           updated_at DATETIME DEFAULT (datetime('now')),
                           deleted_at DATETIME,
                           first_name TEXT NOT NULL,
                           last_name TEXT NOT NULL,
                           address TEXT NOT NULL,
                           email TEXT NOT NULL UNIQUE
);

CREATE TABLE orders (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        created_at DATETIME DEFAULT (datetime('now')),
                        updated_at DATETIME DEFAULT (datetime('now')),
                        deleted_at DATETIME,
                        type TEXT NOT NULL,
                        size TEXT NOT NULL,
                        quantity INTEGER NOT NULL,
                        customer_id INTEGER NOT NULL,
                        status TEXT NOT NULL DEFAULT 'new'
);

CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);

-- Триггер для автоматического обновления updated_at при изменении записи,
-- но вообще обычно всякую логику не рекомендуют выносить на уровень СУБД
CREATE TRIGGER on_update_orders_updated_at
    AFTER UPDATE ON orders
    FOR EACH ROW WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE orders SET updated_at = datetime('now') WHERE id = NEW.id;
END;
