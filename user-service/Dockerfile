FROM ubuntu:latest
LABEL authors="oksas"

ENTRYPOINT ["top", "-b"]


# Gunakan base image Golang
FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o user-service ./cmd/main.go

FROM gcr.io/distroless/base

WORKDIR /app

EXPOSE 8081

COPY --from=builder /app/user-service .

# Jalankan aplikasi
CMD ["/app/user-service"]