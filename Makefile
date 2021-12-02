BIN := "./bin/abf"
DOCKER_IMG="abf:dev"

build:
	go build -v -o $(BIN) ./cmd/abf

run: build
	$(BIN) --config ./configs/config.json

build-img:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -race -count 100 -cover ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.43.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img test lint
