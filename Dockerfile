FROM golang:1.23.0 AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go

CMD ["./main"]