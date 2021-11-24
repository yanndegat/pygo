##
# Pygo
#
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PYGO_GEN_FILES?=$$(find tests -type d -name 'pygo' |grep -v vendor)

all: help

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

unittest: tests/mygolib.so ## run unit tests
	@python3 -m unittest discover tests "*.py"


dist: ## build dist package
	@python3 setup.py sdist bdist_wheel

dist-check: dist ## check if dist is correct
	@twine check dist/*

tests/mygolib.go:
tests/mygolib.so: tests/mygolib.go ## build go shared library for tests
	 @go generate tests/mygolib.go
#@go build -o $@ -buildmode=c-shared tests/mygolib.go

bin:
	@mkdir -p $@

.PHONY: pygo
pygo: go-fmtcheck bin
	@go build -o bin/$@ main.go

pygo-clean:
	@rm -Rf $(PYGO_GEN_FILES)

go-test: go-fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

go-fmt:
	gofmt -w $(GOFMT_FILES)

go-fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

# end
