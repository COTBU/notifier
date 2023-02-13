FROM golang:alpine AS build_base
RUN set -ex && apk update && apk add -q  \
  git unzip build-base autoconf libtool

WORKDIR $GOPATH/src/SOTBI/telegram-notifier

COPY . .

RUN go mod tidy && \
  go build -a -installsuffix cgo -o notifer . && \
  mv notifer /notifer

# Start fresh from a smaller image
FROM alpine:latest

WORKDIR /root/
COPY --from=build_base notifer .
ENTRYPOINT ["./notifer","-config", "/config/notifier.yaml"]