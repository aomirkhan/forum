CREATE TABLE IF NOT EXISTS users (
    name VARCHAR(30) NOT NULL PRIMARY KEY,
    email VARCHAR(45) NOT NULL,
    password VARCHAR(45) NOT NULL
);
CREATE TABLE IF NOT EXISTS cookies (
    Name VARCHAR(30) NOT NULL PRIMARY KEY,
    Id VARCHAR(100) NOT NULL
);