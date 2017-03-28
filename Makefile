EXTENSION_DIR ?= ./ext
APP_NAME = rubberdoc
GOBIN = $(GOPATH)/bin
GLIDE_VERSION := $(shell glide -v 2> /dev/null)

all: install

.PHONY: submodules
submodules:
	git submodule update --init --recursive

.PHONY: drafter
drafter:
	$(MAKE) -C $(EXTENSION_DIR)/drafter

.PHONY: install-glide
install-glide:
ifndef GLIDE_VERSION
	$(shell mkdir -p $(GOBIN))
	$(shell curl https://glide.sh/get | sh)
endif

.PHONY: glide-install
glide-install: install-glide
	$(shell glide install)

.PHONY: go-test
go-test:
	go test $(shell glide novendor) -v

.PHONY: go-gen
go-gen:
	go generate main.go

.PHONY: go-build
go-build:
	go build -o $(GOBIN)/$(APP_NAME) .

.PHONY: go-install
go-install:
	go build -i -o $(GOBIN)/$(APP_NAME) .

.PHONY: clean
clean:
	$(RM) $(GOBIN)/$(APP_NAME)

dep: submodules drafter glide-install
build: clean dep go-gen go-build
install: clean dep go-gen go-install
test: clean dep go-gen go-test
