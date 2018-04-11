# Package configuration
PROJECT = client-go
LIBUAST_VERSION ?= 1.9.1
GOPATH ?= $(shell go env GOPATH)

TOOLS_FOLDER = tools

ifneq ($(OS),Windows_NT)
COPY = cp
else
COPY = copy
endif

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --depth 1 $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

GOGET ?= go get

clean: clean-libuast
clean-libuast:
	find ./  -regex '.*\.[h,c]c?' ! -name 'bindings.h' -exec rm -f {} +

dependencies: cgo-dependencies
ifneq ($(OS),Windows_NT)
cgo-dependencies:
	curl -SL https://github.com/bblfsh/libuast/releases/download/v$(LIBUAST_VERSION)/libuast-v$(LIBUAST_VERSION).tar.gz | tar xz
	mv libuast-v$(LIBUAST_VERSION)/src/* $(TOOLS_FOLDER)/.
	rm -rf libuast-v$(LIBUAST_VERSION)
	$(GOGET) ./...
else
binaries.win64.mingw\lib:
	go get -v github.com/mholt/archiver/cmd/archiver
	cd $(TOOLS_FOLDER) && \
	curl -SLko binaries.win64.mingw.zip https://github.com/bblfsh/libuast/releases/download/v$(LIBUAST_VERSION)/binaries.win64.mingw.zip && \
	$(GOPATH)\bin\archiver open binaries.win64.mingw.zip && \
	del /q binaries.win64.mingw.zip && echo done

cgo-dependencies: binaries.win64.mingw\lib
	$(GOGET) ./...
endif  # !Windows_NT

# $(DEPENDENCIES) it's allowed to file since the code is not compilable
# without libuast.
.IGNORE: $(DEPENDENCIES)
