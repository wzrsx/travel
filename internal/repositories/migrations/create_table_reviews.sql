CREATE TABLE reviews (
    id_review SERIAL PRIMARY KEY,
    description VARCHAR(200),
    estimation INTEGER,
    id_photo INTEGER,
    CONSTRAINT fk_reviews_photos
        FOREIGN KEY (id_photo) REFERENCES photos(id_photo)
);