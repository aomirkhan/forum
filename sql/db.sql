CREATE TABLE users(
    Name VARCHAR(30) NOT NULL,
    Email VARCHAR(50) NOT NULL,
    Password VARCHAR(40) NOT NULL

);
CREATE TABLE posts(
    Post VARCHAR(500) NOT NULL,
    Namae VARCHAR(50) NOT NULL,
    Category VARCHAR(40) NOT NULL,
    Id INT NOT NULL 
);
CREATE TABLE IF NOT EXISTS cookies   ( 
							Id VARCHAR(50),
							lame VARCHAR(50)
                            )



