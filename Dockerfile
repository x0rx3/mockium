FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mockium ./cmd/app.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/mockium /app/mockium

EXPOSE 5000

ENTRYPOINT ["/app/mockium"]
CMD ["-template", "templates"]