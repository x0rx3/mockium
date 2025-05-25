FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mockium ./cmd/app.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/mockium /app/mockium
COPY templates /app/templates

EXPOSE 5000

CMD ["/app/mockium", "-address", "127.0.0.1:5000", "-template", "/app/templates"]