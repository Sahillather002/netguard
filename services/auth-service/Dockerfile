FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-service .
EXPOSE 8081
CMD ["./auth-service"]

# Updated: 2025-07-29T11:30:00

# feat(docker): add Dockerfile for Auth Service - 2025-07-29

# feat(docker): add Dockerfile for Auth Service - 2025-07-29
