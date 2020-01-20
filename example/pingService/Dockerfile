ARG alpine_version=3.11
FROM golang:1.13.6-alpine$alpine_version AS build

RUN apk add --update --no-cache git
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /app/


WORKDIR /app
# Build the gateway

RUN go build -o pingService

FROM alpine:$alpine_version
WORKDIR /app
COPY --from=build /app/pingService /app/
EXPOSE 9600

CMD ["./pingService"]
