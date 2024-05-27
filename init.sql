-- CREATE DATABASE IF NOT EXISTS mydb;
-- USE mydb;

-- create table user  (
--     id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
--     mail VARCHAR(100) NOT NULL,
--     username VARCHAR(100) UNIQUE,
--     pswd VARCHAR(100) NOT NULL, 
-- );

CREATE USER platops WITH PASSWORD 'platops';
CREATE DATABASE platopsdb;
GRANT ALL PRIVILEGES ON DATABASE platopsdb TO platops;