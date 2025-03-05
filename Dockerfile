# Stage 1: Build the Go application
FROM golang:1.24-alpine3.21 AS builder

ARG COMMIT
ARG VERSION

# Set the Current Working Directory inside the container
WORKDIR /build

COPY go.* ./ 
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

COPY . ./
RUN go build \
    -ldflags=" \
    -X gitea.ctplt.ru/catapulto/prostor-sms-balance-exporter/pkg/info.CommitSHA=${COMMIT} \
    -X gitea.ctplt.ru/catapulto/prostor-sms-balance-exporter/pkg/info.Version=${VERSION} \
    " \
    -o app ./cmd 

# Stage 2: Create a small image with the built binary
FROM alpine:3.21.3 AS app

RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
COPY --from=builder /build/app /usr/local/bin/app 
# Copy the Pre-built binary file from the previous stage

WORKDIR /usr/local/bin/
# Command to run the executable
CMD ["./app"]