GOPATH=$(shell go env GOPATH)
# GOLANGCI_LINT_VERSION is the latest version of golangci-lint
# adjusted to match the Go version used in this project
GOLANGCI_LINT_VERSION=v1.61.0

.PHONY: lint
lint: ## run all the lint tools, install golangci-lint if not exist
ifeq (,$(wildcard $(GOPATH)/bin/golangci-lint))
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) > /dev/null
	$(GOPATH)/bin/golangci-lint run --timeout 180s || exit 0
else
	$(GOPATH)/bin/golangci-lint run --timeout 180s || exit 0
endif

.PHONY: vet
vet: ## Field Alignment
ifeq (,$(wildcard $(GOPATH)/bin/fieldalignment))
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go vet -vettool=$(GOPATH)/bin/fieldalignment ./... || exit 0
else
	go vet -vettool=$(GOPATH)/bin/fieldalignment ./... || exit 0
endif

.PHONY: vet-fix
vet-fix: ##If fixed, the annotation for struct fields will be removed
ifeq (,$(wildcard $(GOPATH)/bin/fieldalignment))
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	$(GOPATH)/bin/fieldalignment -fix ./... || exit 0
else
	$(GOPATH)/bin/fieldalignment -fix ./... || exit 0
endif