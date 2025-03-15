# Use a lightweight Go image
FROM golang:1.24.1 as builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o auth-service

# Use a minimal image for the final container
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service .
CMD ["./auth-service"]
