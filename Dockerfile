FROM golang:1.24.3-alpine3.21 AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/myapp ./slog

FROM alpine:3.21.3

COPY --from=builder /app/myapp /usr/local/bin/myapp

CMD ["myapp"]