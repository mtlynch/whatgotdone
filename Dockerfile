FROM node:20.6.1-alpine AS frontend_builder

COPY ./frontend /app/frontend
WORKDIR /app/frontend

ARG NPM_BUILD_MODE="dev"
RUN echo "npm build mode: ${NPM_BUILD_MODE}"

RUN npm install
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.23.1 AS backend_builder

COPY ./backend /app/backend
COPY ./dev-scripts/build-backend /app/dev-scripts/build-backend

WORKDIR /app

ARG GO_BUILD_MODE="dev"

RUN echo "Go build mode: ${GO_BUILD_MODE}"
RUN ./dev-scripts/build-backend "${GO_BUILD_MODE}"

FROM debian:stable-20211011-slim AS litestream_downloader

ARG litestream_version="v0.3.7"
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
COPY --from=backend_builder /app/bin/whatgotdone /app/bin/whatgotdone
COPY --from=litestream_downloader /litestream/litestream /app/litestream
COPY ./creds /app/creds
COPY ./litestream.yml /etc/litestream.yml
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]
