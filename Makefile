BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
APP=skyfka
VERSION=1.0.1

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

default: build

.PHONY:all
all: clean release

.PHONY:install
install:
	go install ./...

.PHONY: build
build:
	go build $(LDFLAGS) -o ./bin/$(APP) ./cmd/$(APP)/...


.PHONY:release
release:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES),\
	$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o ./bin/$(APP)_$(GOOS)_$(GOARCH) ./cmd/$(APP)/...)))
	chmod +x ./bin/*

.PHONY:clean
clean:
	@rm ./bin/*
