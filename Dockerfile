FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем API
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Собираем SEED
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed


FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/api .
COPY --from=builder /app/seed .

EXPOSE 8080
CMD ["./api"]
