CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) UNIQUE NOT NULL,
    JWT VARCHAR(255) NULL,
    authToken VARCHAR(255) NULL,
    authStatus BOOLEAN DEFAULT FALSE,
	avatar	TEXT,
	notify BOOLEAN DEFAULT TRUE
);

CREATE TABLE stickers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    image_path TEXT
);

CREATE TABLE images (
	image_path TEXT,
	id SERIAL PRIMARY KEY,
	userId bigint,
	uploadTime timestamp with time zone,
	like_count INTEGER DEFAULT 0,
	comment_count INTEGER DEFAULT 0
);
