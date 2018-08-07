HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-values/
HAS_GLIDE := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_PLUGIN_DIR)
	cp $(BUILD)/values $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o $(BUILD)/values -ldflags $(LDFLAGS)

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/values -ldflags $(LDFLAGS)
	tar -zcvf $(DIST)/helm-values-linux-$(VERSION).tgz $(BUILD)/values README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/values -ldflags $(LDFLAGS)
	tar -zcvf $(DIST)/helm-values-macos-$(VERSION).tgz $(BUILD)/values README.md LICENSE plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/values.exe -ldflags $(LDFLAGS)
	tar -llzcvf $(DIST)/helm-values-windows-$(VERSION).tgz $(BUILD)/values.exe README.md LICENSE plugin.yaml

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
