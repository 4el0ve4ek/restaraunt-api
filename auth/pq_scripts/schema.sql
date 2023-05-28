CREATE TABLE "user"(
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('customer', 'chef', 'manager')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE session (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    session_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
-- TODO: use session_token index

INSERT INTO "user"(username, email, password_hash, role)
VALUES ('admin', 'admin@mail.ru', 'secret_hash', 'manager');

INSERT INTO "session"(user_id, session_token, expires_at)
VALUES (1, '20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2', '2025-10-10 10:10:10');


-- for setting current timestamp on update
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_customer_modtime BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

----