--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

-- Started on 2025-02-26 21:54:47

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 16478)
-- Name: photos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.photos (
    id_photo integer NOT NULL,
    path_to_photo character varying(200) NOT NULL
);


ALTER TABLE public.photos OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16477)
-- Name: photos_id_photo_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.photos_id_photo_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.photos_id_photo_seq OWNER TO postgres;

--
-- TOC entry 4831 (class 0 OID 0)
-- Dependencies: 217
-- Name: photos_id_photo_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.photos_id_photo_seq OWNED BY public.photos.id_photo;


--
-- TOC entry 220 (class 1259 OID 16485)
-- Name: reviews; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reviews (
    id_review integer NOT NULL,
    description character varying(200),
    estimation integer,
    id_photo integer
);


ALTER TABLE public.reviews OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16484)
-- Name: reviews_id_review_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reviews_id_review_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.reviews_id_review_seq OWNER TO postgres;

--
-- TOC entry 4832 (class 0 OID 0)
-- Dependencies: 219
-- Name: reviews_id_review_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reviews_id_review_seq OWNED BY public.reviews.id_review;


--
-- TOC entry 222 (class 1259 OID 16497)
-- Name: routes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.routes (
    id_route integer NOT NULL,
    name character varying(20),
    yandex_route character varying(200) NOT NULL,
    estimation integer,
    id_review integer
);


ALTER TABLE public.routes OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16496)
-- Name: routes_id_route_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.routes_id_route_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.routes_id_route_seq OWNER TO postgres;

--
-- TOC entry 4833 (class 0 OID 0)
-- Dependencies: 221
-- Name: routes_id_route_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.routes_id_route_seq OWNED BY public.routes.id_route;


--
-- TOC entry 224 (class 1259 OID 16509)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(20) NOT NULL,
    password character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    id_route integer
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 16508)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 4834 (class 0 OID 0)
-- Dependencies: 223
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4656 (class 2604 OID 16481)
-- Name: photos id_photo; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.photos ALTER COLUMN id_photo SET DEFAULT nextval('public.photos_id_photo_seq'::regclass);


--
-- TOC entry 4657 (class 2604 OID 16488)
-- Name: reviews id_review; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews ALTER COLUMN id_review SET DEFAULT nextval('public.reviews_id_review_seq'::regclass);


--
-- TOC entry 4658 (class 2604 OID 16500)
-- Name: routes id_route; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.routes ALTER COLUMN id_route SET DEFAULT nextval('public.routes_id_route_seq'::regclass);


--
-- TOC entry 4659 (class 2604 OID 16512)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 4819 (class 0 OID 16478)
-- Dependencies: 218
-- Data for Name: photos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.photos (id_photo, path_to_photo) FROM stdin;
\.


--
-- TOC entry 4821 (class 0 OID 16485)
-- Dependencies: 220
-- Data for Name: reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reviews (id_review, description, estimation, id_photo) FROM stdin;
\.


--
-- TOC entry 4823 (class 0 OID 16497)
-- Dependencies: 222
-- Data for Name: routes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.routes (id_route, name, yandex_route, estimation, id_review) FROM stdin;
\.


--
-- TOC entry 4825 (class 0 OID 16509)
-- Dependencies: 224
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, email, id_route) FROM stdin;
\.


--
-- TOC entry 4835 (class 0 OID 0)
-- Dependencies: 217
-- Name: photos_id_photo_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.photos_id_photo_seq', 1, false);


--
-- TOC entry 4836 (class 0 OID 0)
-- Dependencies: 219
-- Name: reviews_id_review_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reviews_id_review_seq', 1, false);


--
-- TOC entry 4837 (class 0 OID 0)
-- Dependencies: 221
-- Name: routes_id_route_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.routes_id_route_seq', 1, false);


--
-- TOC entry 4838 (class 0 OID 0)
-- Dependencies: 223
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- TOC entry 4661 (class 2606 OID 16483)
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id_photo);


--
-- TOC entry 4663 (class 2606 OID 16490)
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id_review);


--
-- TOC entry 4665 (class 2606 OID 16502)
-- Name: routes routes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.routes
    ADD CONSTRAINT routes_pkey PRIMARY KEY (id_route);


--
-- TOC entry 4667 (class 2606 OID 16516)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4669 (class 2606 OID 16514)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4670 (class 2606 OID 16491)
-- Name: reviews fk_reviews_photos; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT fk_reviews_photos FOREIGN KEY (id_photo) REFERENCES public.photos(id_photo);


--
-- TOC entry 4671 (class 2606 OID 16503)
-- Name: routes fk_routes_review; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.routes
    ADD CONSTRAINT fk_routes_review FOREIGN KEY (id_review) REFERENCES public.reviews(id_review);


--
-- TOC entry 4672 (class 2606 OID 16517)
-- Name: users fk_users_routes; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_routes FOREIGN KEY (id_route) REFERENCES public.routes(id_route);


-- Completed on 2025-02-26 21:54:47

--
-- PostgreSQL database dump complete
--

