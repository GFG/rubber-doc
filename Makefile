EXTENSION_DIR ?= ./ext
PACKAGE_NAME = rubberdoc
BINARY_DEST = $(GOPATH)/bin/$(PACKAGE_NAME)

all: install

.PHONY: submodules
submodules:
	git submodule update --init --recursive

.PHONY: drafter
drafter:
	cd $(EXTENSION_DIR)/drafter
	$(MAKE) -C $(EXTENSION_DIR)/drafter

.PHONY: glide
glide:
    exists = $(shell glide -v)
    ifndef exists
        $(shell curl https://glide.sh/get | sh)
    endif

.PHONY: glide-install
glide-install:
	$(shell glide install)

.PHONY: go-test
go-test:
	go test $(shell glide novendor) -v

.PHONY: go-gen
go-gen:
	go generate main.go

.PHONY: go-build
go-build:
	go build -o $(BINARY_DEST) .

.PHONY: go-install
go-install:
	go build -i -o $(BINARY_DEST) .

.PHONY: clean
clean:
	$(RM) $(BINARY_DEST)

dep: submodules drafter glide glide-install
build: clean dep go-gen go-build
install: clean dep go-gen go-install
test: dep go-gen go-test
