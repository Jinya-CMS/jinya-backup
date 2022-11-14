FROM quay.imanuel.dev/dockerhub/library---golang:1.19-alpine AS build

# Resolve app dependencies.
WORKDIR /app
COPY . .
RUN apk add git

RUN go build -o /app/jinya-backup jinya-backup

FROM quay.imanuel.dev/dockerhub/library---node:latest AS build-frontend
COPY web/ /app/web
RUN cd /app/web && npm install

FROM quay.imanuel.dev/dockerhub/library---alpine:latest
WORKDIR /app
COPY --from=build /app/jinya-backup /app/jinya-backup
COPY --from=build-frontend /app/web /app/web

# Start server.
ENTRYPOINT ["/app/jinya-backup", "serve"]