FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gomock ./cmd/gomock.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/gomock /app/gomock
COPY templates/ /app/templates/

EXPOSE 8080

ENTRYPOINT ["/app/gomock"]