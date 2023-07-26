CREATE TABLE users(
    Name VARCHAR(30) NOT NULL,
    Email VARCHAR(50) NOT NULL,
    Password VARCHAR(40) NOT NULL

);
CREATE TABLE posts(
    Title VARCHAR(50) NOT NULL,
    Post VARCHAR(500) NOT NULL,
    Namae VARCHAR(50) NOT NULL,
    Category VARCHAR(40) NOT NULL,
    Id INT NOT NULL 
);
CREATE TABLE IF NOT EXISTS cookies   ( 
	Id VARCHAR(50),
	lame VARCHAR(50)
);
CREATE TABLE comments(
    Name VARCHAR(30) NOT NULL,
    Text VARCHAR(200) NOT NULL,
    Id VARCHAR(40) NOT NULL,
    Comid INT
);
CREATE TABLE comlikes(
    Name VARCHAR(30) NOT NULL,
    Comid VARCHAR(3) NOT NULL,
     Id INT

);
CREATE TABLE comdislikes(
    Name VARCHAR(30) NOT NULL,
    Comid VARCHAR(3) NOT NULL,
     Id INT

);
CREATE TABLE likes(
    Name VARCHAR(30) NOT NULL,
    Postid VARCHAR(3) NOT NULL
    -- FOREIGN KEY (Postid) REFERENCES posts(Id) 

);
CREATE TABLE dislikes(
    Name VARCHAR(30) NOT NULL,
    Postid VARCHAR(3) NOT NULL
    -- FOREIGN KEY (Postid) REFERENCES posts(Id) 

);