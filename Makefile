# go option
GO             ?= go
SHELL          := /usr/bin/env bash
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
PKG            := $(GO) mod download
# TODO(jaredallard): infer from Git tag
APP_VERSION    := $(shell git describe --match 'v[0-9]*' --tags --abbrev=0 --always HEAD)
LDFLAGS        := -w -s -X github.com/tritonmedia/pkg/app.Version=$(APP_VERSION)
GOFLAGS        :=
GO_EXTRA_FLAGS := -v
TAGS           :=
BINDIR         := $(CURDIR)/bin
PKGDIR         := github.com/tritonmedia/api
CGO_ENABLED    := 1
BENCH_FLAGS    := "-bench=Bench $(BENCH_FLAGS)"
TEST_TAGS      ?= tm_teste
LOG            := "$(CURDIR)/scripts/make-log-wrapper.sh"

.PHONY: default
default: build

.PHONY: version
version:
	@echo "$(APP_VERSION)"

.PHONY: release
release:
	@$(LOG) info "Building official release"
	./scripts/gobin.sh github.com/goreleaser/goreleaser

.PHONY: pre-commit
pre-commit: fmt

.PHONY: build
build: gogenerate gobuild

.PHONY: test
test:
	./scripts/test.sh

.PHONY: dep
dep:
	@$(LOG) info "Installing dependencies using '$(PKG)'"
	$(PKG)

.PHONY: gogenerate
gogenerate:
	@$(LOG) info "Running 'go generate'"
	$(GO) generate ./...

.PHONY: gobuild
gobuild:
	@$(LOG) info "Building releases into ./bin"
	@mkdir -p $(BINDIR)
	CGO_ENABLED=$(CGO_ENABLED) "$(GO)" build -o "$(BINDIR)/" -ldflags "$(LDFLAGS)" $(GO_EXTRA_FLAGS) "$(PKGDIR)/..."

.PHONY: docker-build-init
docker-build-init:
	docker buildx create --use
	docker run --rm --privileged docker/binfmt:66f9012c56a8316f9244ffd7622d7c21c1f6f28d

.PHONY: docker-build-push
docker-build-push:
	@$(LOG) info "Building and push docker image"
	docker buildx build --platform amd64,arm64 -t "tritonmedia/api:$(APP_VERSION)" --push .

.PHONY: fmt
fmt:
	@$(LOG) info "Running goimports"
	find  . -path ./vendor -prune -o -type f -name '*.go' -print | xargs ./scripts/gobin.sh golang.org/x/tools/cmd/goimports -w
	@$(LOG) info "Running shfmt"
	./scripts/gobin.sh mvdan.cc/sh/v3/cmd/shfmt -l -w -s .

