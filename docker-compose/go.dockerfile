FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache wget

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
