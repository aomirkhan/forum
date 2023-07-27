FROM golang:latest
LABEL maintainer="Rizabek Zhampeisov, Azamat Omirhan, rizabekzhampeisov440@gmail.com"
LABEL description ="docker for forum"
WORKDIR /app

COPY ./ ./
RUN apt-get update && apt-get -y install sudo
RUN go mod download
RUN sudo apt-get -y update
RUN sudo apt-get -y upgrade
RUN sudo apt-get install -y sqlite3 libsqlite3-dev

RUN sqlite3 sql/database.db < sql/db.sql
# RUN go get github.com/gofrs/uuid
# RUN go get github.com/mattn/go-sqlite3
# RUN go get golang.org/x/crypto
RUN go build -o forum .
# EXPOSE 8000
CMD ["./forum"]
# RUN go build -o ascii-art-web-dockerize .