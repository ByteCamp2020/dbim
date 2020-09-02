export GO111MODULE := on
GOOS := $(if $(GOOS),$(GOOS),"")
GOARCH := $(if $(GOARCH),$(GOARCH),"")
GOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
CGOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO     := $(GOENV) go
CGO    := $(CGOENV) go
GOTEST := TEST_USE_EXISTING_CLUSTER=false NO_PROXY="${NO_PROXY},testhost" go test
SHELL    := /usr/bin/env bash
COMMIT_SHA=$(shell git rev-parse --short HEAD)


fmt:
	$(CGOENV) go fmt ./...
vet:
	$(CGOENV) go vet ./...
# build:
# 	$(CGOENV) go build -o bin/logic ./src/cmd/logic/main.go
# 	$(CGOENV) go build -o bin/worker ./src/cmd/worker/main.go
# 	$(CGOENV) go build -o bin/comet ./src/cmd/comet/main.go
# 	$(CGOENV) go build -o bin/dbworker ./src/cmd/dbworker/main.go
build:
    GOARCH=amd64 GOOS=linux go build -o dep/logic/main ./src/cmd/logic/main.go
	GOARCH=amd64 GOOS=linux go build -o dep/worker/main ./src/cmd/worker/main.go
	GOARCH=amd64 GOOS=linux go build -o dep/comet/main ./src/cmd/comet/main.go
# 	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o dep/dbworker/dbworker ./src/cmd/dbworker/main.go
image:
	docker build -t bin/logic:${COMMIT_SHA} .
