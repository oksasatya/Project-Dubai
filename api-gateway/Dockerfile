FROM ubuntu:latest
LABEL authors="oksas"

ENTRYPOINT ["top", "-b"]


# Gunakan base image Golang
FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o api-gateway ./cmd/main.go

FROM gcr.io/distroless/base

WORKDIR /app

EXPOSE 8081

COPY --from=builder /app/api-gateway .

# Jalankan aplikasi
CMD ["/app/api-gateway"]