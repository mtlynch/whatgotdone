FROM node:12.18.4-alpine AS frontend_builder

COPY ./frontend /app/frontend
WORKDIR /app/frontend

ARG NPM_BUILD_MODE="development"
RUN npm install
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.16.7 AS backend_builder

COPY ./backend /app/backend

WORKDIR /app

ARG GO_BUILD_TAGS="dev"
RUN go build --tags "$GO_BUILD_TAGS" -o /app/main backend/main.go

FROM golang:1.16.7

COPY --from=frontend_builder /app/frontend/dist /app/frontend/dist
COPY --from=backend_builder /app/main /app/main

WORKDIR /app

ENTRYPOINT ["/app/main"]
