FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the entire project for monorepo build
COPY . .

# Build the application from the service directory
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/api-gateway/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 