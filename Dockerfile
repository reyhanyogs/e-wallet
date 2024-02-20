# Build stage
FROM golang:1.21.6-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY main.env .
COPY start.sh .
COPY db/migrations ./db/migrations
COPY wait-for .
RUN chmod a+x start.sh
RUN chmod a+x wait-for
RUN mkdir logs

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]