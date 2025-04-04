BEGIN;


ALTER TABLE IF EXISTS public.photos DROP CONSTRAINT IF EXISTS fk_photos_routes;

ALTER TABLE IF EXISTS public.reviews DROP CONSTRAINT IF EXISTS fk_reviews_routes;

ALTER TABLE IF EXISTS public.reviews DROP CONSTRAINT IF EXISTS fk_user;

ALTER TABLE IF EXISTS public.routes DROP CONSTRAINT IF EXISTS fk_routes_users;



DROP TABLE IF EXISTS public.photos;

CREATE TABLE IF NOT EXISTS public.photos
(
    id_photo serial NOT NULL,
    path_to_photo character varying(200) COLLATE pg_catalog."default",
    id_route integer,
    CONSTRAINT photos_pkey PRIMARY KEY (id_photo)
);

DROP TABLE IF EXISTS public.reviews;

CREATE TABLE IF NOT EXISTS public.reviews
(
    id_review serial NOT NULL,
    description text COLLATE pg_catalog."default",
    estimation integer,
    date_review timestamp without time zone,
    id_route integer,
    id_user integer NOT NULL,
    CONSTRAINT reviews_pkey PRIMARY KEY (id_review),
    CONSTRAINT reviews_id_user_key UNIQUE (id_user)
);

DROP TABLE IF EXISTS public.routes;

CREATE TABLE IF NOT EXISTS public.routes
(
    id_route serial NOT NULL,
    route_name character varying(100) COLLATE pg_catalog."default",
    route_place character varying(120) COLLATE pg_catalog."default",
    route_description text COLLATE pg_catalog."default",
    yandex_route character varying(500) COLLATE pg_catalog."default" NOT NULL,
    estimation real,
    path_to_photo_preview character varying(200) COLLATE pg_catalog."default",
    id_user integer,
    CONSTRAINT routes_pkey PRIMARY KEY (id_route)
);

DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id_user serial NOT NULL,
    username character varying(20) COLLATE pg_catalog."default" NOT NULL,
    password character varying(50) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id_user),
    CONSTRAINT unq_user UNIQUE (id_user),
    CONSTRAINT users_email_key UNIQUE (email)
);

ALTER TABLE IF EXISTS public.photos
    ADD CONSTRAINT fk_photos_routes FOREIGN KEY (id_route)
    REFERENCES public.routes (id_route) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.reviews
    ADD CONSTRAINT fk_reviews_routes FOREIGN KEY (id_route)
    REFERENCES public.routes (id_route) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.reviews
    ADD CONSTRAINT fk_user FOREIGN KEY (id_user)
    REFERENCES public.users (id_user) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;
CREATE INDEX IF NOT EXISTS reviews_id_user_key
    ON public.reviews(id_user);


ALTER TABLE IF EXISTS public.routes
    ADD CONSTRAINT fk_routes_users FOREIGN KEY (id_user)
    REFERENCES public.users (id_user) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

-- Вставка данных в таблицу users
INSERT INTO users (username, email, password) VALUES 
('vasiliy', 'vasiliy.rybakov.06@mail.ru', 'admin'),
('vitaliy', 'vitaliy@gmail.com', '20.06.2006'),
('sergey', 'serzh.rybakov.06@mail.ru', 'admin'),
('alex', 'alex@example.com', 'password123'),
('maria', 'maria@example.com', 'securepass');

-- Вставка данных в таблицу routes
INSERT INTO routes (route_name, route_place, route_description, yandex_route, estimation, path_to_photo_preview, id_user) VALUES 
('Москва-Пасad', 'Москва', 'Красивый маршрут через центр', 'yandex.ru/maps/123', 4.5, 'routes/moscow-pasad/preview.jpg', 1),
('Золотое кольцо', 'Ярославль', 'Тур по древним городам', 'yandex.ru/maps/456', 4.8, 'routes/golden-ring/preview.jpg', 2),
('Байкал', 'Иркутск', 'Путешествие вокруг озера', 'yandex.ru/maps/789', 4.9, 'routes/baikal/preview.jpg', 3),
('Сочи-Красная поляна', 'Сочи', 'Горный маршрут', 'yandex.ru/maps/101', 4.3, 'routes/sochi/preview.jpg', 4),
('Казань', 'Казань', 'Исторический центр', 'yandex.ru/maps/112', 4.6, 'routes/kazan/preview.jpg', 5);

-- Вставка данных в таблицу reviews (исправлено - используется id_user вместо username)
INSERT INTO reviews (description, estimation, date_review, id_route, id_user) VALUES 
('Отличный маршрут, рекомендую!', 5, NOW(), 1, 2),
('Очень красивые места', 4, NOW(), 1, 3),
('Прекрасные виды на озеро', 5, NOW(), 3, 4),
('Интересная экскурсия', 4, NOW(), 4, 5),
('Много исторических мест', 5, NOW(), 5, 1);

-- Вставка данных в таблицу photos
INSERT INTO photos (path_to_photo, id_route) VALUES 
('routes/moscow-pasad/photo1.jpg', 1),
('routes/moscow-pasAD/photo2.jpg', 1),
('routes/golden-ring/photo1.jpg', 2),
('routes/baikal/photo1.jpg', 3),
('routes/baikal/photo2.jpg', 3),
('routes/sochi/photo1.jpg', 4),
('routes/kazan/photo1.jpg', 5),
('routes/kazan/photo2.jpg', 5);

COMMIT;