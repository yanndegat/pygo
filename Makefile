##
# Pygo
#
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PYGO_GEN_FILES?=$$(find tests -type d -name 'pygo' |grep -v vendor)
PYGO_LOG?=ERROR

export PATH := $(shell pwd)/bin:$(PATH)

all: help

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

unittest: go-test tests/simplelib.so tests/mylibgo/pygo/mygolib.so ## run unit tests
	@echo "--- RUNNING PYTHON UNIT TESTS  ---"
	@python3 -m unittest discover tests "*.py"

dist: ## build dist package
	@python3 setup.py sdist bdist_wheel

dist-check: dist ## check if dist is correct
	@twine check dist/*

tests/simplelib.go:
tests/simplelib.so: tests/simplelib.go ## build simple go shared library for tests
	 @go build -o $@ -buildmode=c-shared tests/simplelib.go

tests/mylibgo/mygolib.go:
tests/mylibgo/pygo/mygolib.so: bin/pygo tests/mylibgo/mygolib.go ## build go shared library for tests
	@PYGO_LOG=$(PYGO_LOG) go generate generate ./...

bin: ## create ./bin directory
	@mkdir -p $@

.PHONY: pygo
pygo: go-fmtcheck bin ## build pygo
	@go build -o bin/$@ main.go

bin/pygo: pygo ## build ./bin/pygo

pygo-clean: ## clean auto generated pygo files
	@rm -Rf $(PYGO_GEN_FILES)

go-test: go-fmtcheck ## run unit tests on go project
	@echo "--- RUNNING GO TESTS  ---"
	@go test -i $(TEST) || exit 1
	@echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

go-fmt: ## formats all go source files
	gofmt -w $(GOFMT_FILES)

go-fmtcheck: ## checks that go source files are properly formatted
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"


# end
