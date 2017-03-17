EXTENSION_DIR ?= ./ext

submodules:
	git submodule update --init --recursive

drafter:
	cd $(EXTENSION_DIR)/drafter
	$(MAKE) -C $(EXTENSION_DIR)/drafter

glide:
    exists = $(shell glide -v)
    ifndef exists
         $(shell curl https://glide.sh/get | sh)
    endif

glide-install: glide $(shell glide install)

dep: submodules drafter glide-install

go-test:
	go test $(shell glide novendor) -v && go test $(shell glide novendor) -v ./...

go-gen:
	go generate main.go

go-build:
	go build -o rubberdoc .

go-install:
	go install ./...

build: dep go-gen go-build
install: dep go-gen go-install
test: dep go-gen go-test

.PHONY: build test install