SHELL:=/bin/bash
TEST?=$$(go list ./azuredevops/internal/acceptancetests |grep -v 'vendor')
UNITTEST?=$$(go list ./... |grep -v 'vendor')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=azuredevops
TESTTIMEOUT=180m
TESTTAGS=all

ifeq ($(GOPATH),)
	GOPATH:=$(shell go env GOPATH)
endif

.EXPORT_ALL_VARIABLES:
  TF_SCHEMA_PANIC_ON_ERROR=1

default: build

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	@echo "GOPATH: $(GOPATH)"
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	go install github.com/bflad/tfproviderdocs@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GOPATH)/bin" v1.27.0

build: fmtcheck check-vendor-vs-mod
	go install

fmt:
	@echo "==> Fixing source code with gofmt..."
	@echo "# This logic should match the search logic in scripts/gofmtcheck.sh"
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@if command -v golangci-lint; then (golangci-lint run ./...); else ($(GOPATH)/bin/golangci-lint run ./...); fi

test: fmtcheck
	go test -tags "all" -i $(UNITTEST) || exit 1
	echo $(UNITTEST) | \
    		xargs -t -n4 go test -tags "all" $(TESTARGS) -timeout=60s -parallel=4

testacc: fmtcheck
	@echo "==> Sourcing .env file if avaliable"
	if [ -f .env ]; then set -o allexport; . ./.env; set +o allexport; fi; \
	TF_ACC=1 go test -tags "$(TESTTAGS)" $(TEST) -v $(TESTARGS) -timeout 120m

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

install:
	./scripts/build.sh --SkipTests --Install

check-vendor-vs-mod: ## Check that go modules and vendored code are on par
	@echo "==> Checking that go modules and vendored dependencies match..."
	go mod vendor
	@if [[ `git status --porcelain vendor` ]]; then \
		echo "ERROR: vendor dir is not on par with go modules definition." && \
		exit 1; \
	fi

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

ci: check-vendor-vs-mod lint test

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck lint tools test-compile website website-lint website-test
