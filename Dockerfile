# syntax=docker/dockerfile:1

FROM golang:1.24

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./api ./cmd/gabcgen

EXPOSE 8080 

CMD ["/src/api"]
