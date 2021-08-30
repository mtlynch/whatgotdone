FROM node:12.18.4-alpine AS frontend_builder

COPY ./frontend /app/frontend
WORKDIR /app/frontend

ARG NPM_BUILD_MODE="development"
RUN npm install --only=production
RUN npm run build -- --mode "$NPM_BUILD_MODE"

FROM golang:1.16.7-alpine

COPY --from=frontend_builder /app/frontend/dist /app/frontend/dist
COPY ./backend /app/backend
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

ARG GO_BUILD_TAGS="dev"
RUN go build --tags "$GO_BUILD_TAGS" -o /app/main backend/main.go

ENTRYPOINT ["/app/main"]
