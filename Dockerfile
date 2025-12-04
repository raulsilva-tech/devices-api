# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build the api
RUN CGO_ENABLED=0 GOOS=linux go build -o devices-api ./cmd/api

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/devices-api .

EXPOSE 8080

CMD ["./devices-api"]
