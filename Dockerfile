# syntax=docker/dockerfile:1.0-experimental
FROM golang:1.14-alpine as builder
ARG VERSION
WORKDIR /src

RUN apk add --no-cache make

# Add go.mod and go.sum first to maximize caching
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# Build our application
RUN make build APP_VERSION=${VERSION}

FROM alpine:3.12
ENTRYPOINT ["/usr/bin/api"]

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

COPY --from=builder /src/bin/api /usr/bin/api
