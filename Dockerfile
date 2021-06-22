FROM golang:1.16-alpine AS build
WORKDIR /build
COPY . .
RUN go get -v ./...
RUN go build -v -o guildapi cmd/guildapi/main.go

FROM alpine:latest AS final
WORKDIR /app
COPY --from=build /build/guildapi .
EXPOSE 80
ENV BINDADDRESS=0.0.0.0:80
ENTRYPOINT ["/app/guildapi"]