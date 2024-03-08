#Build stage
FROM golang:1.22.0-alpine3.18 AS builder
WORKDIR /app  
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz


# Run stage 
FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .

COPY db/migrations ./migrations
COPY start.sh .

RUN chmod +x /app/start.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]

