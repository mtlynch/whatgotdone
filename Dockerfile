FROM node:10.15.3 AS frontend_builder

COPY ./web/frontend /app/web/frontend
WORKDIR /app/web/frontend
RUN npm run build

FROM golang:1.11.8

COPY --from=frontend_builder /app/web/frontend/dist /app/web/frontend/dist

EXPOSE $PORT
ENTRYPOINT go run web/main.go