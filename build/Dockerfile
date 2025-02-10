FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../. .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/avito_merch ./cmd


FROM debian:stable-slim

WORKDIR /app

COPY ./cmd/.env /app

COPY --from=build-stage /app/bin/avito_merch bin

ENTRYPOINT ["./bin"]