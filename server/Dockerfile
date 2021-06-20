FROM dart:stable AS build

# Resolve app dependencies.
WORKDIR /app
COPY pubspec.* ./
RUN dart pub get

# Copy app source code and AOT compile it.
COPY . .
# Ensure packages are still up-to-date if anything has changed
RUN dart pub get --offline
RUN dart compile exe bin/server.dart -o bin/server
RUN dart compile exe bin/console.dart -o bin/console

FROM alpine
COPY --from=build /runtime/ /
COPY --from=build /app/bin/server /bin/
COPY --from=build /app/bin/console /bin/
COPY --from=build /app/docker/entrypoint.sh /

# Start server.
ENTRYPOINT "/entrypoint.sh"