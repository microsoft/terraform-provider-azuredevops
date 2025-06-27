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
  GO111MODULE=on

default: build

tools:
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	go install github.com/bflad/tfproviderdocs@latest
	go install github.com/katbyte/terrafmt@latest
	go install mvdan.cc/gofumpt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GOPATH)/bin" v1.64.8

build: fmtcheck depscheck
	go install

fmt:
	@echo "==> Fixing source code with gofmt..."
	@echo "# This logic should match the search logic in scripts/gofmtcheck.sh"
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

fumpt:
	@echo "==> Fixing source code with Gofumpt..."
	# This logic should match the search logic in scripts/gofmtcheck.sh
	find . -name '*.go' | grep -v vendor | xargs gofumpt -s -w

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

terrafmt:
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@if command -v terrafmt; \
		then (find azuredevops | egrep "_test.go" | sort | while read f; do terrafmt fmt -f $$f; done) \
		else (find azuredevops | egrep "_test.go" | sort | while read f; do $(GOPATH)/bin/terrafmt fmt -f $$f; done); \
	  fi
	@echo "==> Fixing website terraform blocks code with terrafmt..."
	@if command -v terrafmt; \
		then (find . | egrep html.markdown | sort | while read f; do terrafmt fmt $$f; done); \
		else (find . | egrep html.markdown | sort | while read f; do $(GOPATH)/bin/terrafmt fmt $$f; done); \
	  fi

terrafmt-check:
	./scripts/terrafmt.sh

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

test: fmtcheck
	go test -v ./...

testacc: fmtcheck
	@echo "==> Sourcing .env file if available"
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

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

ci: depscheck lint test website-lint

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/
	@echo "==> Checking documentation for errors..."
	@tfproviderdocs check -provider-name=azuredevops
#-require-resource-subcategory \
#		-allowed-resource-subcategories-file website/allowed-subcategories

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

scaffold-website:
	./scripts/scaffold-website.sh

.PHONY: build test testacc vet fmt fmtcheck lint tools test-compile website website-lint website-test
