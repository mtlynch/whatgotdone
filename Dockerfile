FROM node:10.15.3 AS frontend_builder

COPY ./web/frontend /app/web/frontend
WORKDIR /app/web/frontend

ARG NPM_BUILD_MODE="development"
RUN npm install
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.11.8

COPY --from=frontend_builder /app/web/frontend/dist /app/web/frontend/dist
COPY ./auth /app/auth
COPY ./datastore /app/datastore
COPY ./handlers /app/handlers
COPY ./types /app/types
COPY ./web/*.go /app/web
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

ARG GO_BUILD_TAGS="dev"
ARG USE_FRESH_REACTIONS_KEY="0"
RUN set -xe && \
  PER_USER_REACTIONS_KEY="perUserReactions" && \
  if test "$USE_FRESH_REACTIONS_KEY" = "1"; then PER_USER_REACTIONS_KEY="${PER_USER_REACTIONS_KEY}-$(date +'%s')"; fi && \
  go build \
  --tags "$GO_BUILD_TAGS" \
  -ldflags "-X github.com/mtlynch/whatgotdone/datastore.perUserReactionsKey=${PER_USER_REACTIONS_KEY}" \
  -o /app/main \
  web/main.go

ENTRYPOINT /app/main