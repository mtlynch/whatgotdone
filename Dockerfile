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

FROM golang:1.24.8 AS litestream_builder

RUN apt-get update && apt-get install -y git

WORKDIR /app

RUN git clone https://github.com/mtlynch/litestream.git && \
    cd litestream && \
    git checkout 489d05e8c92d014b0475604efae3771135b9151e && \
    CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/litestream

FROM alpine:3.15

RUN apk add --no-cache bash

COPY --from=frontend_builder /app/frontend/dist /app/frontend/dist
COPY --from=backend_builder /app/bin/whatgotdone /app/bin/whatgotdone
COPY --from=litestream_builder /app/litestream/litestream /app/litestream
COPY ./creds /app/creds
COPY ./litestream.yml /etc/litestream.yml
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]
