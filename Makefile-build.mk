VALIDATE_DEPS     = github.com/golang/lint/golint
TEST_DEPS         = github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml

build: clean validate test compile

$(INTS):
	@if [ -f $(INTEGRATIONS_DIR)/$@/Makefile ]; then \
		make -C $(INTEGRATIONS_DIR)/$@ $$TARGET ;\
	else \
		echo "=== Main === [ $$TARGET ] - $@: no Makefile found. Skipping." ;\
	fi

validate-deps:
	@echo "=== Main === [ validate-deps ]: installing validation dependencies..."
	@go get -v $(VALIDATE_DEPS)

validate:
	@echo "=== Main === [ validate ]: running validation for: $(INTS)"
ifeq ($(INTEGRATIONS),all)
	@TARGET=validate-only $(MAKE) --no-print-directory validate-deps $(INTS)
else
	@TARGET=validate $(MAKE) --no-print-directory $(INTS)
endif

compile:
	@echo "=== Main === [ compile ]: building the following integrations: $(INTS)"
ifeq ($(INTEGRATIONS),all)
	@TARGET=compile $(MAKE) --no-print-directory $(INTS)
else
	@TARGET=compile $(MAKE) --no-print-directory $(INTS)
endif

test-deps:
	@echo "=== Main === [ test-deps ]: installing testing dependencies..."
	@go get -v $(TEST_DEPS)

test:
	@echo "=== Main === [ test ]: running unit tests for the following integrations: $(INTS)"
ifeq ($(INTEGRATIONS),all)
	@$(MAKE) --no-print-directory test-deps
	gocov test ./integrations/... | gocov-xml > coverage.xml
else
	@TARGET=test $(MAKE) --no-print-directory $(INTS)
endif

.PHONY: build $(INTS) validate-deps validate compile test-deps test
