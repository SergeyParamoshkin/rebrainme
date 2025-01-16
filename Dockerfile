FROM golang:1.23.4 AS build 

WORKDIR /src
COPY *.go ./
RUN go mod downloaded
RUN go build ...


FROM alpine:3.18 AS app 

COPY --from=build /src/... /app
ENTRYPOINT /app