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
    path_to_photo_preview VARCHAR(100),
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

-- Фотки в маршруте (не превью)
CREATE TABLE photos (
    id_photo SERIAL PRIMARY KEY,
    path_to_photo VARCHAR(200),
    id_route INTEGER,
    CONSTRAINT fk_photos_routes FOREIGN KEY (id_route) REFERENCES routes(id_route)
);
-- Users
INSERT INTO users (username, email, password) VALUES ('vasiliy', 'vasiliy.rybakov.06@mail.ru', 'admin');
INSERT INTO users (username, email, password) VALUES ('vitaliy', 'vitaliy@gmail.com', '20.06.2006');
INSERT INTO users (username, email, password) VALUES ('sergey', 'serzh.rybakov.06@mail.ru', 'admin');

-- Routes
INSERT INTO routes (route_name, route_place, route_description, yandex_route, estimation, path_to_photo_preview, id_user) VALUES ('Москва-Пасад','Москва-Пасад','Москва-Пасад', 'yandex.ru', 4, 'sadf/asdad/asd/asd/', 1);
INSERT INTO routes (route_name, route_place, route_description, yandex_route, estimation, path_to_photo_preview, id_user) VALUES ('Москва-Пасад','Москва-Пасад','Москва-Пасад', 'yandex.ru', 4, 'sadf/asdad/asd/asd/', 2);
INSERT INTO routes (route_name, route_place, route_description, yandex_route, estimation, path_to_photo_preview, id_user) VALUES ('Москва-Пасад','Москва-Пасад','Москва-Пасад', 'yandex.ru', 4, 'sadf/asdad/asd/asd/', 3);

-- Reviews
INSERT INTO reviews (description, estimation, id_route) VALUES ('класс', 7, 1);
INSERT INTO reviews (description, estimation, id_route) VALUES ('супер', 9, 1);
INSERT INTO reviews (description, estimation, id_route) VALUES ('кайф', 10, 3);