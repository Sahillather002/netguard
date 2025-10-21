FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]

# Updated: 2025-07-29T10:00:00

# feat(docker): add Dockerfile for API Gateway - 2025-07-29

# feat(docker): add Dockerfile for API Gateway - 2025-07-29
