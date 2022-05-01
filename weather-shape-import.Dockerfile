FROM golang:1.18-alpine as build

RUN mkdir -p /app/src & mkdir -p /app/bin
WORKDIR /app/src

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/bin/weather-zone-import /app/src/cmd/weather-zone-import/weather-zone-import.go


FROM alpine AS production

WORKDIR /app

RUN mkdir -p /app/shapefiles
COPY ./shapefiles/*.zip /app/shapefiles/
COPY ./deploy/scripts/import-zones.sh /app/bin/import-zones.sh
RUN chmod 755 /app/bin/import-zones.sh
COPY --from=build /app/bin/weather-zone-import /app/bin/weather-zone-import

ENV MONGO_HOST host.docker.internal
ENV MONGO_PORT 27017

ENTRYPOINT ["/bin/sh", "/app/bin/import-zones.sh"]