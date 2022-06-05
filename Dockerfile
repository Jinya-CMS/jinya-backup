FROM quay.imanuel.dev/dockerhub/library---dart:sdk AS build

# Resolve app dependencies.
WORKDIR /app
COPY pubspec.* ./
RUN dart pub get

# Copy app source code and AOT compile it.
COPY server .
# Ensure packages are still up-to-date if anything has changed
RUN dart pub get --offline
RUN dart compile exe bin/server.dart -o bin/server
RUN dart compile exe bin/console.dart -o bin/console

FROM quay.imanuel.dev/dockerhub/library---node:latest AS build-frontend
COPY server/frontend /app/frontend
RUN cd /app/frontend && npm install

FROM quay.imanuel.dev/dockerhub/library---alpine:latest
COPY --from=build /runtime/ /
COPY --from=build /app/bin/server /bin/
COPY --from=build /app/bin/console /bin/
COPY --from=build /app/docker/entrypoint.sh /
COPY --from=build-frontend /app/frontend /frontend

# Start server.
ENTRYPOINT "/entrypoint.sh"