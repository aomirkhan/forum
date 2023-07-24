FROM golang:latest
LABEL maintainer="Rizabek Zhampeisov, Azamat Omirhan, rizabekzhampeisov440@gmail.com"
LABEL description ="docker for forum"
WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go get github.com/gofrs/uuid
RUN go get github.com/mattn/go-sqlite3
RUN go get golang.org/x/crypto
# RUN go build -o ascii-art-web-dockerize .