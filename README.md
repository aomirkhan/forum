# FORUM

## TABLE OF CONTENTS
* General Info
* Usage
* Usage with Docker
* Contributors


### General Info
* The aim of this project is to construct a simple forum where users are able to to add posts, write comments in those created posts and add likes/dislikes in the posts/comments.

* Details of how code works can be found in .go files in the form of comments.

### Usage
* To run the app:
``` 
cd forum
``` 
sqlite3 sql/database.db < sql/db.sql
```
go run .
```
### Usage with Docker
* To build the image you need to write: 
``` 
cd forum
``` 
docker build -t example .
```
Afterwards, to run the container you need to write:
``` 
docker run -d -p 8000:8000 example
```
After that just follow this link:
```
http://localhost:8000
```
### Contributors
* @rzhampeis 
* @aomirhan  