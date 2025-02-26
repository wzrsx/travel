Create table users(
	id SERIAL PRIMARY KEY,
	username VARCHAR(20) NOT NULL,
	password VARCHAR(50) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE
)