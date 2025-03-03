-- Установка кодировки и стандартных строк
SET client_encoding = 'UTF8';
SET standard_conforming_strings = 'on';

-- Установка пути поиска
SELECT pg_catalog.set_config('search_path', '', false);

-- Создание базы данных
CREATE DATABASE "webRoutes" 
WITH TEMPLATE = template0 
ENCODING = 'UTF8' 
LOCALE_PROVIDER = libc 
LOCALE = 'ru-RU';

-- Подключение к базе данных
\c webRoutes

-- Создание таблицы users
CREATE TABLE public.users (
    id_user integer NOT NULL,
    username character varying(20) NOT NULL,
    password character varying(50) NOT NULL,
    email character varying(100) NOT NULL
);

-- Создание последовательности для users
CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Привязка последовательности к столбцу id_user
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id_user;

-- Установка значения по умолчанию для id_user
ALTER TABLE ONLY public.users ALTER COLUMN id_user SET DEFAULT nextval('public.users_id_seq'::regclass);

-- Добавление первичного ключа для users
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id_user);

-- Добавление уникального ограничения для email
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

-- Создание таблицы routes
CREATE TABLE public.routes (
    id_route integer NOT NULL,
    name character varying(20),
    yandex_route character varying(200) NOT NULL,
    estimation integer,
    id_user integer
);

-- Создание последовательности для routes
CREATE SEQUENCE public.routes_id_route_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Привязка последовательности к столбцу id_route
ALTER SEQUENCE public.routes_id_route_seq OWNED BY public.routes.id_route;

-- Установка значения по умолчанию для id_route
ALTER TABLE ONLY public.routes ALTER COLUMN id_route SET DEFAULT nextval('public.routes_id_route_seq'::regclass);

-- Добавление первичного ключа для routes
ALTER TABLE ONLY public.routes
    ADD CONSTRAINT routes_pkey PRIMARY KEY (id_route);

-- Добавление внешнего ключа для routes
ALTER TABLE ONLY public.routes
    ADD CONSTRAINT fk_routes_users FOREIGN KEY (id_user) REFERENCES public.users(id_user);

-- Создание таблицы reviews
CREATE TABLE public.reviews (
    id_review integer NOT NULL,
    description character varying(200),
    estimation integer,
    id_route integer
);

-- Создание последовательности для reviews
CREATE SEQUENCE public.reviews_id_review_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Привязка последовательности к столбцу id_review
ALTER SEQUENCE public.reviews_id_review_seq OWNED BY public.reviews.id_review;

-- Установка значения по умолчанию для id_review
ALTER TABLE ONLY public.reviews ALTER COLUMN id_review SET DEFAULT nextval('public.reviews_id_review_seq'::regclass);

-- Добавление первичного ключа для reviews
ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id_review);

-- Добавление внешнего ключа для reviews
ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_reviews_routes FOREIGN KEY (id_route) REFERENCES public.routes(id_route);

-- Установка значений последовательностей
SELECT pg_catalog.setval('public.users_id_seq', 1, true);
SELECT pg_catalog.setval('public.routes_id_route_seq', 1, false);
SELECT pg_catalog.setval('public.reviews_id_review_seq', 1, false);