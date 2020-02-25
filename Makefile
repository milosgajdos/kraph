BUILD=go build
CLEAN=go clean
INSTALL=go install
BUILDPATH=./_build
PACKAGES=$(shell go list ./... | grep -v /examples/)

kraphctl: builddir
	go build -o "$(BUILDPATH)/kraphctl" "cmd/kraphctl/main.go"

all: dep check test kraphctl

builddir:
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

.PHONY: clean kraphctl
