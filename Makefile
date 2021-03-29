PROJECTNAME = $(shell basename "$(PWD)")

include project.properties

PB_REL     = https://github.com/protocolbuffers/protobuf/releases
PB_DIR     = ~/.protoc

BASE_DIR   = $(shell pwd)
BIN_DIR    = $(BASE_DIR)/.bin
TMP_DIR    = $(BASE_DIR)/.tmp
STUB_DIR   = $(BASE_DIR)/pkg/pb
MOCK_DIR   = $(BASE_DIR)/internal/mocks
ASSET_DIR  = $(BASE_DIR)/assets
PROTO_DIR  = $(ASSET_DIR)/proto

PROTOFILES = $(shell find $(ASSET_DIR) -type f -name '*.proto')

GO_PKGS    = $(shell go list ./...)
GO_FILES   = $(shell find . -type f -name '*.go')
GO_LDFLAGS = "-s -w -X main.App=$(APP) -X main.Version=$(VERSION)"

MOCKFILES  = \
  $(GOPATH)/src/github.com/hashicorp/raft/fsm.go \

TEST_OPTS  = -timeout 300s -cover -coverprofile=$(COVERAGE_FILE) -failfast

TEST_PKGS  = \
  ./internal/cluster         \
  ./internal/gateway         \
  ./internal/util            \
  ./pkg/client               \

BINARY_FILE   = $(BIN_DIR)/$(APP)
COVERAGE_FILE = $(BIN_DIR)/test-coverage.out

BIN_MARKER   = $(TMP_DIR)/bin.marker
MOCKS_MARKER = $(TMP_DIR)/mocks.marker
STUBS_MARKER = $(TMP_DIR)/stubs.marker
TEST_MARKER  = $(TMP_DIR)/test.marker
VET_MARKER   = $(TMP_DIR)/vet.marker

ifndef VERBOSE
.SILENT:
endif

.PHONY: all help clean deps

all: help

## help: Prints help
help: Makefile
	@echo "Choose a command in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## test: Runs unit tests.
test: $(TEST_MARKER)

$(TEST_MARKER): mocks $(GO_FILES)
	@echo "Executing unit tests..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(BIN_DIR); go test $(TEST_OPTS) $(TEST_PKGS)
	@mkdir -p $(@D)
	@touch $@

## vet: Runs Go vet
vet: $(VET_MARKER)

$(VET_MARKER): $(GO_FILES)
	@echo "Running Go vet..."
	@GOBIN=$(BIN_DIR); go vet $(GO_PKGS)
	@gosec ./...
	@mkdir -p $(@D)
	@touch $@

## build: Builds the binaries
build: $(BIN_MARKER)

$(BIN_MARKER): $(GO_FILES) stubs test
	@echo "Building the binaries..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(BIN_DIR) \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -tags=$(tags) -ldflags=$(GO_LDFLAGS) -o $(BINARY_FILE) cmd/*.go
	@mkdir -p $(@D)
	@touch $@

## clean: Cleans the project
clean:
	@echo "Cleaning the project..."
	@rm -rf $(BIN_DIR) $(TMP_DIR)

## stubs: Generates stubs
stubs: $(STUBS_MARKER)

$(STUBS_MARKER): $(PROTOFILES)
	@echo "Generating Protobuf stubs..."
	@mkdir -p $(STUB_DIR) $(TMP_DIR)
	@$(PB_DIR)/bin/protoc \
      --proto_path=$(PROTO_DIR)/$(PROTO_VERSION) \
      --proto_path=${GOPATH}/src \
      --proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      --go_out=$(STUB_DIR) \
      --go-grpc_out=$(STUB_DIR) \
      --validate_out=lang=go:$(STUB_DIR) \
      --grpc-gateway_out=logtostderr=true:$(STUB_DIR) \
      --swagger_out=logtostderr=true:$(STUB_DIR) \
     $?
	@mkdir -p $(@D)
	@touch $@

## mocks: Generates mocks
mocks: $(MOCKS_MARKER)

$(MOCKS_MARKER): $(MOCKFILES)
	echo "Generating mocks..."
	mockgen -source $(GOPATH)/src/github.com/hashicorp/raft/fsm.go \
        -destination $(MOCK_DIR)/github.com/hashicorp/raft/fsm.go
	@mkdir -p $(@D)
	@touch $@

## image: Builds the Docker image
image: build
	@echo "Building Docker image..."
	$(eval tmp := $(shell mktemp -d))
	@cp -r $(BINARY_FILE) $(tmp)
	@cp -r $(ASSET_DIR)/docker/Dockerfile $(tmp)
	@docker build --build-arg VERSION=$(VERSION) \
      -t $(DOCKER_NS)/$(APP):$(VERSION) $(tmp)

## bench: Runs benchmarks
bench: $(GO_FILES)
	@echo "Executing benchmarks..."
	@GOBIN=$(BIN_DIR); go test -bench=. $(TEST_PKGS)

## deps: Installs dependencies
deps:
	@mkdir -p $(PB_DIR)
	@curl -s -LO $(PB_REL)/download/v$(PB_VER)/protoc-$(PB_VER)-linux-x86_64.zip
	@unzip -o -qq protoc-$(PB_VER)-linux-x86_64.zip -d $(PB_DIR)
	@rm protoc-$(PB_VER)-linux-x86_64.zip
	@go mod tidy
	@go install github.com/golang/mock/mockgen@v1.5.0
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/envoyproxy/protoc-gen-validate
	@go install \
      github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
      github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
      google.golang.org/protobuf/cmd/protoc-gen-go \
      google.golang.org/grpc/cmd/protoc-gen-go-grpc
