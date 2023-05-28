CREATE TABLE dish (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    is_available BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "order" (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    special_requests TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
--     FOREIGN KEY (user_id) REFERENCES users(id) -- это живет в другом сервисе, мы ничего не знаем про табличку users
);

CREATE TABLE order_dish (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    dish_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES "order"(id),
    FOREIGN KEY (dish_id) REFERENCES "dish"(id) ON DELETE CASCADE
);


-- обновляет столбик updated_at  на каждое изменение строки
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_dish_modtime BEFORE UPDATE ON "dish" FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();
CREATE TRIGGER update_order_modtime BEFORE UPDATE ON "order" FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();
