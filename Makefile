LIBUAST_VERSION=0.2.0

GO_CMD = go
GO_BUILD = $(GO_CMD) get
GO_CLEAN = $(GO_CMD) clean
GO_TEST = $(GO_CMD) test -v

.PHONY: all clean deps build test

all: build

clean:
	find ./ -name '*.[h,c]' ! -name 'bindings.h' -exec rm -f {} +
	$(GO_CLEAN)

deps:
	curl -SL https://github.com/bblfsh/libuast/releases/download/v$(LIBUAST_VERSION)/libuast-v$(LIBUAST_VERSION).tar.gz | tar xz
	mv libuast-v$(LIBUAST_VERSION)/src/* .
	rm -rf libuast-v$(LIBUAST_VERSION)

build: deps
	$(GO_BUILD) ./...

test: build
	$(GO_TEST) ./...
