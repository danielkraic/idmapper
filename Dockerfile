# stage 1
FROM golang:1.12-alpine as build
RUN apk --no-cache add build-base git make
WORKDIR /src
COPY . /src
RUN make build

# stage 2
FROM alpine:3
WORKDIR /app
COPY --from=build /src/idmapperapp /app/
COPY --from=build /src/config-example.yaml /app/config.yaml
EXPOSE 80
ENTRYPOINT ./idmapperapp