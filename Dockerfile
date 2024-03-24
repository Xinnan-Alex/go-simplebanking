# Build Stage
FROM golang:1.22.1-alpine3.19 AS buidler

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=buidler /app/main .

EXPOSE 8080
CMD ["/app/main"]