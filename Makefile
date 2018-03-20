SRCS=$(shell find . -type f -regextype posix-egrep ! -regex "^\./vendor/.*" -regex ".*\.go$$")
PKGS=$(shell go list ./...)

scoreserver: $(SRCS)
	go build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	(! gofmt -s -d $(SRCS) | grep ^)
	go vet ./...
	golint -set_exit_status $(PKGS)
	go test -race ./...
