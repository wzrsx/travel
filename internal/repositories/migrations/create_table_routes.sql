Create TABLE routes(
	id_route SERIAL PRIMARY KEY,
	name VARCHAR(20),
	yandex_route VARCHAR(100) NOT NULL,
	estimation INTEGER,
	id_review INTEGER,
	CONSTRAINT fk_routes_review
        FOREIGN KEY (id_review) REFERENCES reviews(id_review)
)