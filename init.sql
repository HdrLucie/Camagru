CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    JWT VARCHAR(255) NULL,
    authToken VARCHAR(255) NULL,
    authStatus BOOLEAN
);