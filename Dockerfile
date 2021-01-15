# syntax=docker/dockerfile:1.0-experimental
FROM golang:1.15-alpine as builder
ARG VERSION
ENV GOCACHE "/go-build-cache"
WORKDIR /src

# Why: Pinning versions in a docker container isn't
# really worth the effort.
# hadolint ignore=DL3018
RUN apk add --no-cache make bash git libc-dev gcc

COPY . .

# --mount here allows us to cache the packages even if
# it's invalidated by go.mod,go.sum
RUN --mount=type=cache,target=/go/pkg make dep

# --mount here allows us to save go build cache across builds
# but also needed to use the package cache above
RUN --mount=type=cache,target=/go-build-cache --mount=type=cache,target=/go/pkg make build APP_VERSION=${VERSION}

FROM alpine:3.13
ENTRYPOINT ["/usr/bin/api"]
CMD ["server"]

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

COPY --from=builder /src/bin/api /usr/bin/api
