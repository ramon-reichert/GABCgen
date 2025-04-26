# syntax=docker/dockerfile:1

FROM golang:1.24

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o ./api ./cmd

EXPOSE 8080 

CMD ["/src/api"]
