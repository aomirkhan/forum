CREATE TABLE likes(
    Name VARCHAR(30) NOT NULL,
    Postid VARCHAR(3) NOT NULL,
    -- FOREIGN KEY (Postid) REFERENCES posts(Id) 

);
CREATE TABLE dislikes(
    Name VARCHAR(30) NOT NULL,
    Postid VARCHAR(3) NOT NULL,
    -- FOREIGN KEY (Postid) REFERENCES posts(Id) 

);