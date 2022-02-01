GOLANGCI_VERSION = 1.44.0
GOLANGCI = vendor/golangci-lint/$(GOLANGCI_VERSION)/golangci-lint

install-golangci:
	test -s $(GOLANGCI) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $(GOLANGCI)) v$(GOLANGCI_VERSION)

lint: install-golangci
	$(GOLANGCI) run
