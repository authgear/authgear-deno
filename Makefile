GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)
IMAGE ?= quay.io/theauthgear/authgear-deno:$(GIT_HASH)

.PHONY: vendor
vendor:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2
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
	docker build --pull --file ./cmd/server/Dockerfile --tag $(IMAGE) .

.PHONY: push-image
push-image:
	docker push $(IMAGE)
