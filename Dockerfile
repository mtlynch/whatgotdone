FROM node:12.18.4-alpine AS frontend_builder

COPY ./frontend /app/frontend
WORKDIR /app/frontend

ARG NPM_BUILD_MODE="development"
RUN npm install
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.16.7 AS backend_builder

COPY ./backend /app/backend

WORKDIR /app/backend

ARG GO_BUILD_TAGS="dev"

RUN GOOS=linux GOARCH=amd64 \
    go build \
      -tags "netgo $GO_BUILD_TAGS" \
      -ldflags '-w -extldflags "-static"' \
      -o /app/main \
      main.go

FROM debian:stable-20211011-slim AS litestream_downloader

ARG litestream_version="v0.3.6"
ARG litestream_binary_tgz_filename="litestream-${litestream_version}-linux-amd64-static.tar.gz"

WORKDIR /litestream

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget
RUN wget "https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/${litestream_binary_tgz_filename}"
RUN tar -xvzf "${litestream_binary_tgz_filename}"

FROM alpine:3.15

RUN apk add --no-cache bash

COPY --from=frontend_builder /app/frontend/dist /app/frontend/dist
COPY --from=backend_builder /app/main /app/main
COPY --from=litestream_downloader /litestream/litestream /app/litestream
COPY ./litestream.yml /etc/litestream.yml
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]
