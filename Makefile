BINDIR := $(PWD)/bin
PATH:=$(PWD)/bin:$(PATH)

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test -race -tags='testing' ./...

.PHONY: build
build: generate
	go build -o build/ethereum-proxy ./cmd/main.go

.PHONY: run
run: build
	build/ethereum-proxy -c config/local

.PHONY: clean
clean:
	rm build/ethereum-proxy

.PHONY: lint
lint:
	golangci-lint run --build-tags="testing" ./...

HAS_golangci-lint := $(shell command -v golangci-lint || command -v $(BINDIR)/golangci-lint;)
HAS_gomock := $(shell command -v mockgen || command -v $(BINDIR)/mockgen;)

.PHONY: docker-up
docker-up:
	docker-compose up

.PHONY: bootstrap
bootstrap:
ifndef HAS_golangci-lint
	@echo "Installing golangci-lint"
	GOBIN=$(BINDIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint
endif
ifndef HAS_gomock
	@echo "Installing gomock"
	GOBIN=$(BINDIR) go install github.com/golang/mock/mockgen
endif

