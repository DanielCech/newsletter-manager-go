include ./Makefile.vars
GO := $(shell which go)

.PHONY:
	run  \
	test \
	build

all: fmt vet build

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

run: RUN_ARGS=--help
run: fmt vet
	$(GO) run ./cmd/api $(RUN_ARGS)

test: generate lint
	$(GO) test ./... -cover

lint: generate
	golangci-lint run

init-api: CMD=""
init-api:
ifeq ($(CMD),rest)
	cp cmd/api/rest/main.go cmd/api/main.go
	rm -rf cmd/api/rest
	rm -rf cmd/api/graphql
	rm -rf api/graphql
	rm -rf domain/user/postgres/dataloader
else ifeq ($(CMD),graphql)
	cp cmd/api/graphql/main.go cmd/api/main.go
	rm -rf cmd/api/graphql
	rm -rf cmd/api/rest
	rm -rf api/rest
else
	$(error CMD is not set)
endif

build: BUILD_OUTPUT=./bin/api
build: generate
	CGO_ENABLED=0 $(GO) build -ldflags "-X main.version=$(APP_VERSION)" -o $(BUILD_OUTPUT) ./cmd/api

generate:
	$(GO) generate ./...

start-local:
	@echo Storage is running at port: $(STORAGE_PORT)
	@echo ...
	docker-compose -f docker-compose.yml up

build-image:
ifdef CONTAINER_REGISTRY
	@echo Building container image...
	docker build -t $(CONTAINER_REGISTRY)/newsletter-manager-go:$(IMAGE_TAG) --build-arg APP_VERSION=$(APP_VERSION) .
else
	$(info missing required value CONTAINER_REGISTRY)
endif
