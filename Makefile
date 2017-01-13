EXTENSION_DIR ?= ./ext

build-drafter:
	git submodule update --init --recursive
	cd $(EXTENSION_DIR)/drafter
	$(MAKE) -C $(EXTENSION_DIR)/drafter

install: build-drafter

