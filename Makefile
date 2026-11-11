# Development tools use the published standard without runtime dependencies.
include standard.mk

CHECK_EXTRA := build integration-compile

.PHONY: build integration-compile integration
build:
	$(GO) build ./...

integration-compile:
	$(GO) test -tags=integration -run '^$$' ./...

# The script owns an isolated disposable Zabbix/PostgreSQL stack.
integration:
	bash scripts/integration.sh -v
