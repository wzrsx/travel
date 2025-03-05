-- Таблица users
CREATE TABLE users (
    id_user SERIAL PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    password VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

-- Таблица routes
CREATE TABLE routes (
    id_route SERIAL PRIMARY KEY,
    route_name VARCHAR(20),
    route_place VARCHAR(30),
    route_description VARCHAR(200),
    yandex_route VARCHAR(200) NOT NULL,
    estimation INTEGER,
    path_to_photo VARCHAR(100),
    id_user INTEGER,
    CONSTRAINT fk_routes_users FOREIGN KEY (id_user) REFERENCES users(id_user)
);

-- Таблица reviews
CREATE TABLE reviews (
    id_review SERIAL PRIMARY KEY,
    description VARCHAR(200),
    estimation INTEGER,
    id_route INTEGER,
    CONSTRAINT fk_reviews_routes FOREIGN KEY (id_route) REFERENCES routes(id_route)
);
-- Users
INSERT INTO users (username, email, password) VALUES ('vasiliy', 'vasiliy.rybakov.06@mail.ru', 'admin');
INSERT INTO users (username, email, password) VALUES ('vitaliy', 'vitaliy@gmail.com', '20.06.2006');
INSERT INTO users (username, email, password) VALUES ('sergey', 'serzh.rybakov.06@mail.ru', 'admin');

-- Routes
INSERT INTO routes (route_name, yandex_route, estimation, path_to_photo, id_user) VALUES ('Москва-Пасад', 'yandex.ru', 7, 'sadf/asdad/asd/asd/', 1);
INSERT INTO routes (route_name, yandex_route, estimation, path_to_photo, id_user) VALUES ('Питер-карелия', 'yandex.ru', 9, 'sadf/asdad/asd/asd/', 2);
INSERT INTO routes (route_name, yandex_route, estimation, path_to_photo, id_user) VALUES ('Брюссель', 'yandex.ru', 3, 'sadf/asdad/asd/asd/', 3);

-- Reviews
INSERT INTO reviews (description, estimation, id_route) VALUES ('класс', 7, 1);
INSERT INTO reviews (description, estimation, id_route) VALUES ('супер', 9, 1);
INSERT INTO reviews (description, estimation, id_route) VALUES ('кайф', 10, 3);



ALTER TABLE routes ADD COLUMN last_selected TIMESTAMP DEFAULT CURRENT_TIMESTAMP;



CREATE OR REPLACE FUNCTION update_last_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_selected = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_select_last_updated
BEFORE SELECT ON routes
FOR EACH ROW
EXECUTE FUNCTION update_last_selected();


--Удаление каждые 24 часа
CREATE EXTENSION pg_cron;

SELECT cron.schedule(
    'delete_old_routes', -- Имя задачи
    'EVERY 12 HOURS',    -- Запускать каждые 12 часов
    $$DELETE FROM routes WHERE last_updated < NOW() - INTERVAL '1 day'$$
);