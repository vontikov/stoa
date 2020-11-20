PROJECTNAME = $(shell basename "$(PWD)")

include project.properties

DOCKER_NS  = github.com/vontikov
APP        = stoa-example

BASE_DIR   = $(shell pwd)
BIN_DIR    = $(BASE_DIR)/.bin
BIN_FILE   = $(BIN_DIR)/$(APP)

GO_FILES   = $(shell find . -type f -name '*.go')
GO_LDFLAGS = "-s -w -X main.Version=$(VERSION)"
GO_FILES   = $(shell find . -type f -name '*.go')

clean:
	@rm -r $(BIN_DIR)

build: $(GO_FILES)
	@echo "Building the binaries..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(BIN_DIR) \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -tags=$(tags) -ldflags=$(GO_LDFLAGS) -o $(BIN_FILE) main.go

image: build
	@echo "Building Docker image..."
	$(eval tmp := $(shell mktemp -d))
	@cp -r $(BIN_FILE) $(tmp)
	@cp -r Dockerfile $(tmp)
	@docker build --build-arg VERSION=$(VERSION) \
      -t $(DOCKER_NS)/$(APP):$(VERSION) $(tmp)

