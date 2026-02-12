GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)
BUILD_CONTEXT ::= .
DOCKERFILE ::= ./cmd/server/Dockerfile
IMAGE_TAG ?= quay.io/theauthgear/authgear-deno:$(GIT_HASH)

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.5.0
	go mod download

.PHONY: start
start:
	go run ./cmd/server

.PHONY: build
build:
	go build -o authgear-deno -tags "osusergo netgo static_build timetzdata" ./cmd/server

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./pkg/...

.PHONY: test
test:
	go test ./...

.PHONY: check-tidy
check-tidy:
	$(MAKE) fmt
	go mod tidy
	git status --porcelain | grep '.*'; test $$? -eq 1

.PHONY: build-image
build-image:
	docker build --pull --file $(DOCKERFILE) --tag $(IMAGE_TAG) $(BUILD_CONTEXT)

.PHONY: gh-actions-env
gh-actions-env:
	@printf "BUILD_CONTEXT=%s\nDOCKERFILE=%s\nIMAGE_TAG=%s\n" $(BUILD_CONTEXT) $(DOCKERFILE) $(IMAGE_TAG)
