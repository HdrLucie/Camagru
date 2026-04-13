FROM golang:alpine AS builder

WORKDIR /app

COPY backend/srcs/go.mod backend/srcs/go.sum ./
RUN go mod download

COPY backend/srcs/ .
RUN go build -o serveur .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/serveur .
COPY frontend/srcs/ ./frontend/srcs/

EXPOSE 8080

CMD ["./serveur"]
