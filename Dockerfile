FROM harbor.ulbricht.casa/library/golang:1.21-alpine AS build

WORKDIR /app

COPY . .

RUN apk add git

RUN go build -o /app/jinya-backup jinya-backup

FROM harbor.ulbricht.casa/library/node:latest AS build-frontend

COPY web/ /app/web

RUN cd /app/web && npm install

FROM harbor.ulbricht.casa/library/alpine:latest

WORKDIR /app

COPY --from=build /app/jinya-backup /app/jinya-backup
COPY --from=build-frontend /app/web /app/web

ENTRYPOINT ["/app/jinya-backup", "serve"]