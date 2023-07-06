FROM golang:1.20-alpine3.17 as builder
WORKDIR /build
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz  && tar xvf migrate.linux-amd64.tar.gz
COPY go.* . 
RUN go mod download
COPY . ./
RUN go build -o app ./cmd

FROM alpine:3.17.3 as app
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
COPY --from=builder /build/app /usr/local/bin/app 
COPY --from=builder /build/migrate /usr/local/bin/migrate

WORKDIR /usr/local/bin/

CMD ["app"]
