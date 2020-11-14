FROM mcr.microsoft.com/dotnet/core/sdk:3.1-alpine AS builder
WORKDIR /opt/jinya/
ADD . /opt/jinya/
RUN dotnet build --configuration=Release
RUN apk update
RUN apk add nodejs npm yarn
RUN yarn
RUN cd ./wwwroot && yarn build

FROM mcr.microsoft.com/dotnet/core/aspnet:3.1
WORKDIR /opt/jinya/
COPY --from=builder /opt/jinya/ /opt/jinya/

EXPOSE 5001
EXPOSE 5000

CMD ["/usr/bin/dotnet", "/opt/jinya/bin/Release/netcoreapp3.1/Jinya.Backup.dll"]