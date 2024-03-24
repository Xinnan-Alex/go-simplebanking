# Build Stage
FROM golang:1.22.1-alpine3.19 AS buidler

WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN apk add --no-cache bash
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=buidler /app/main .
COPY --from=buidler /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x /app/start.sh
COPY db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]