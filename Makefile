export GO111MODULE := on
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0
GOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO     := $(GOENV) go
GOTEST := TEST_USE_EXISTING_CLUSTER=false NO_PROXY="${NO_PROXY},testhost" go test
SHELL    := /usr/bin/env bash
COMMIT_SHA=$(shell git rev-parse --short HEAD)


fmt:
	$(CGOENV) go fmt ./...
vet:
	$(CGOENV) go vet ./...
build:
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o dep/logic/logic ./src/cmd/logic/main.go
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o dep/worker/worker ./src/cmd/worker/main.go
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o dep/comet/comet ./src/cmd/comet/main.go
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o dep/dbworker/dbworker ./src/cmd/dbworker/main.go
image:
	docker build -t comet dep/comet
	docker build -t logic dep/logic
	docker build -t worker dep/worker
	docker build -t dbworker dep/dbworker
