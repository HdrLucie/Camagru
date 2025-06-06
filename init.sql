CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    JWT VARCHAR(255) NULL,
    authToken VARCHAR(255) NULL,
    authStatus BOOLEAN,
	avatar	TEXT
);

CREATE TABLE stickers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    image_path TEXT
);

CREATE TABLE images (
	image_path TEXT,
	id SERIAL PRIMARY KEY NOT NULL,
	userId SERIAL PRIMARY KEY NOT NULL,
	uploadTime timestamp with time zone NOT NULL
);
