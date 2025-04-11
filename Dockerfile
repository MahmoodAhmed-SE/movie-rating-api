# Builder stage
FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o movie-api ./cmd/main.go 

# Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/movie-api .
EXPOSE 4040
CMD ["./movie-api"]