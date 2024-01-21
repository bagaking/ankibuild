BUNDLE_RC_FILES := $(shell find . -name "*.file_bundle_rc")
BUNDLE_FILES := $(BUNDLE_RC_FILES:%.file_bundle_rc=%.bundle)

.PHONY: bundle example

example:
	go run .

bundle: $(BUNDLE_FILES)

%.bundle: %.file_bundle_rc
	@file_bundle -v -s -i $<