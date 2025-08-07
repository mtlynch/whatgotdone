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

FROM debian:stable-20240311-slim AS litestream_downloader

ARG TARGETPLATFORM
ARG litestream_version="v0.3.13"

WORKDIR /litestream

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget

RUN set -x && \
    if [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
      ARCH="arm7" ; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
      ARCH="arm64" ; \
    else \
      ARCH="amd64" ; \
    fi && \
    set -u && \
    litestream_binary_tgz_filename="litestream-${litestream_version}-linux-${ARCH}.tar.gz" && \
    wget "https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/${litestream_binary_tgz_filename}" && \
    mv "${litestream_binary_tgz_filename}" litestream.tgz
RUN tar -xvzf litestream.tgz

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
