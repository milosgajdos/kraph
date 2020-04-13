BUILD=go build
CLEAN=go clean
INSTALL=go install
BUILDPATH=./_build
PACKAGES=$(shell go list ./... | grep -v /examples/)

kctl: dir
	go build -o "$(BUILDPATH)/kctl" "cmd/kctl/main.go"

all: dep check test kctl

dir:
	mkdir -p $(BUILDPATH)

clean:
	rm -rf $(BUILDPATH)

dep:
	go get ./...

check:
	go vet ./...

test:
	for pkg in ${PACKAGES}; do \
		go test -coverprofile="../../../$$pkg/coverage.txt" -covermode=atomic $$pkg || exit; \
	done

build:
	go build ./...

.PHONY: clean kctl
