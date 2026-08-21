FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tree-tracker ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=build /tree-tracker /tree-tracker

EXPOSE 8080
ENTRYPOINT ["/tree-tracker"]
