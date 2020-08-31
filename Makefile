export GO111MODULE := on
GOOS := $(if $(GOOS),$(GOOS),"")
GOARCH := $(if $(GOARCH),$(GOARCH),"")
GOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
CGOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO     := $(GOENV) go
CGO    := $(CGOENV) go
GOTEST := TEST_USE_EXISTING_CLUSTER=false NO_PROXY="${NO_PROXY},testhost" go test
SHELL    := /usr/bin/env bash

fmt:
	$(CGOENV) go fmt ./...
build:
	$(CGOENV) go build -o bin/logic ./src/cmd/logic/main.go
	$(CGOENV) go build -o bin/worker ./src/cmd/worker/main.go
	$(CGOENV) go build -o bin/comet ./src/cmd/comet/main.go
