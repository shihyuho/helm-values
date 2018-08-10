HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-values/
HAS_GLIDE := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := values

.PHONY: install
install: bootstrap test build
	mkdir -p $(HELM_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hookInstall
hookInstall: bootstrap test build

.PHONY: test
test:
	go test -v

.PHONY: build
build:
	go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS)

.PHONY: dist
dist:
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	cp README.md $(BUILD) && cp LICENSE $(BUILD) && sed -E 's/(version: )"(.+)"/\1"$(VERSION)"/g' plugin.yaml > $(BUILD)/plugin.yaml
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/helm-$(BINARY)-linux-$(VERSION).tgz $(BINARY) README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -zcvf $(DIST)/helm-$(BINARY)-macos-$(VERSION).tgz $(BINARY) README.md LICENSE plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo
	tar -C $(BUILD) -llzcvf $(DIST)/helm-$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe README.md LICENSE plugin.yaml

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor

.PHONY: clean
clean:
	rm -rf _*
