ARG alpine_version=3.11
FROM golang:1.12.14-alpine$alpine_version AS build

RUN apk add --update --no-cache git
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /app/


WORKDIR /app/cmd/cannon
# Build the gateway

RUN go build -o cannon

FROM alpine:$alpine_version
WORKDIR /app
COPY --from=build /app/cmd/cannon /app/
COPY --from=build /app/configs /app/configs
RUN apk add --no-cache util-linux

EXPOSE 7000

CMD ["./cannon -c /app/configs/config.yaml"]
