PACKAGE_TYPES     ?= deb rpm
PUSH              ?= false
PROJECT_NAME       = newrelic-integrations
BINS_PREFIX        = nr
BINS_DIR           = $(TARGET_DIR)/bin/linux_amd64
SOURCE_DIR         = $(TARGET_DIR)/source
PACKAGE_DIR        = $(TARGET_DIR)/package
DEB_FILENAME      := $(PROJECT_NAME)_$(VERSION)_amd64.deb
RPM_FILENAME      := $(PROJECT_NAME)-$(subst -,_,$(VERSION))-1.x86_64.rpm
FPM_COMMON_OPTIONS = --verbose -C $(SOURCE_DIR) -s dir -n $(PROJECT_NAME) -v $(VERSION) --prefix "" --iteration 1 --license GPL --vendor "New Relic, Inc." -m Jenkins --url https://newrelic.com/infrastructure --config-files /etc/newrelic-infra/ --description "This is a fancy description about a package that does lots of fancy things. I need this description to be longer so I'm just writing that."
FPM_DEB_OPTIONS    = -t deb -p $(PACKAGE_DIR)/deb/$(DEB_FILENAME)
FPM_RPM_OPTIONS    = -t rpm -p $(PACKAGE_DIR)/rpm/$(RPM_FILENAME) --epoch 0 --rpm-summary "This is a fancy summary."

package: create-bins prep-pkg-env $(PACKAGE_TYPES)

create-bins:
	@for I in $(INTS); do \
		if [ $$I != "example" ]; then \
			PACKAGE=$$(go list ./integrations/$$I/... 2>&1) ;\
			if echo $$PACKAGE | grep -Eq ".*matched\ no\ packages$$"; then \
				echo "=== Main === [ create-bins ]: no Go files found for $$I. Skipping." ;\
			else \
				echo "=== Main === [ create-bins ]: creating binary for $$I..." ;\
				go build -v -ldflags '-X main.buildVersion=$(VERSION)' -o $(BINS_DIR)/$(BINS_PREFIX)-$$I $$PACKAGE || exit 1 ;\
			fi ;\
		fi ;\
	done
	@echo ""

prep-pkg-env:
	@if [ ! -d $(BINS_DIR) ]; then \
		echo "=== Main === [ prep-pkg-env ]: no built binaries found. Run 'make create-bins'" ;\
		exit 1 ;\
	fi
	@echo "=== Main === [ prep-pkg-env ]: preparing a clean packaging environment..."
	@rm -rf $(SOURCE_DIR)
	@mkdir -p $(SOURCE_DIR)/var/db/newrelic-infra/newrelic-integrations/bin $(SOURCE_DIR)/etc/newrelic-infra/integrations.d
	@echo "=== Main === [ prep-pkg-env ]: adding built binaries and configuration and definition files..."
	@for BIN in $$(ls $(BINS_DIR)); do \
		cp $(BINS_DIR)/$$BIN $(SOURCE_DIR)/var/db/newrelic-infra/newrelic-integrations/bin ;\
		cp $(INTEGRATIONS_DIR)/$${BIN#$(BINS_PREFIX)-}/*.yml $(SOURCE_DIR)/var/db/newrelic-infra/newrelic-integrations/ ;\
		cp $(INTEGRATIONS_DIR)/$${BIN#$(BINS_PREFIX)-}/*.sample $(SOURCE_DIR)/etc/newrelic-infra/integrations.d/ ;\
	done
	@echo ""

deb: prep-pkg-env
	@echo "=== Main === [ deb ]: building DEB package..."
	@mkdir -p $(PACKAGE_DIR)/deb
	@fpm $(FPM_COMMON_OPTIONS) $(FPM_DEB_OPTIONS) .
	@echo ""
	@PACKAGE_TYPE=deb FILE=$(PACKAGE_DIR)/deb/$(DEB_FILENAME) $(MAKE) --no-print-directory push

rpm: prep-pkg-env
	@echo "=== Main === [ rpm ]: building RPM package..."
	@mkdir -p $(PACKAGE_DIR)/rpm
	@fpm $(FPM_COMMON_OPTIONS) $(FPM_RPM_OPTIONS) .
	@echo ""
	@PACKAGE_TYPE=rpm FILE=$(PACKAGE_DIR)/rpm/$(RPM_FILENAME) $(MAKE) --no-print-directory push

push:
ifeq ($(PUSH),true)
	@echo "=== Main === [ push ]: uploading $$PACKGE_TYPE to PackageCloud..."
	@package_cloud push rtorrents/newrelic-integrations-internal/$$PACKAGE_TYPE $$FILE
else
	@echo "=== Main === [ push ]: upload to PackageCloud disabled. Skipping."
endif
	@echo ""

.PHONY: package create-bins prep-pkg-env deb rpm push
