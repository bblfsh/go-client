# Package configuration
PROJECT = client-go
LIBUAST_VERSION=1.0.0

# Including ci Makefile
MAKEFILE = Makefile.main
CI_REPOSITORY = https://github.com/src-d/ci.git
CI_FOLDER = .ci

$(MAKEFILE):
	@git clone --quiet $(CI_REPOSITORY) $(CI_FOLDER); \
	cp $(CI_FOLDER)/$(MAKEFILE) .;

-include $(MAKEFILE)

clean: clean-libuast
clean-libuast:
	find ./ -name '*.[h,c]' ! -name 'bindings.h' -exec rm -f {} +

dependencies: cgo-dependencies
cgo-dependencies:
	curl -SL https://github.com/bblfsh/libuast/releases/download/v$(LIBUAST_VERSION)/libuast-v$(LIBUAST_VERSION).tar.gz | tar xz
	mv libuast-v$(LIBUAST_VERSION)/src/* .
	rm -rf libuast-v$(LIBUAST_VERSION)
	$(GOGET) .

# $(DEPENDENCIES) it's allowed to file since the code it's not compilable
# without libuast.
.IGNORE: $(DEPENDENCIES)
