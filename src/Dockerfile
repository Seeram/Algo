FROM golang:1.22.2

WORKDIR /app

COPY * ./

RUN go mod download

RUN go build -o /app/algo

EXPOSE 8080