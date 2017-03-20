EXTENSION_DIR ?= ./ext

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
	go test $(shell glide novendor) -v && go test $(shell glide novendor) -v ./...

.PHONY: go-gen
go-gen:
	go generate main.go

.PHONY: go-build
go-build:
	go build -o rubberdoc .

.PHONY: go-install
go-install:
	go install ./...

.PHONY: clean
clean:
	$(RM) rubberdoc

dep: submodules drafter glide glide-install
build: dep go-gen go-build
install: dep go-gen go-install
test: dep go-gen go-test
