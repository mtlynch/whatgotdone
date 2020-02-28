FROM node:10.15.3-alpine AS frontend_builder

COPY ./frontend /app/frontend
WORKDIR /app/frontend

ARG NPM_BUILD_MODE="development"
RUN npm install
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.13.5-buster

COPY --from=frontend_builder /app/frontend/dist /app/frontend/dist
COPY ./backend /app/backend
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./cache/go-modules/ /go/pkg/mod

WORKDIR /app

ARG GO_BUILD_TAGS="dev"
RUN go build --tags "$GO_BUILD_TAGS" -o /app/main backend/main.go

ENTRYPOINT ["/app/main"]