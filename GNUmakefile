TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=cloud66
VERSION=$(shell git describe --tags --always)

default: build

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	go install github.com/bflad/tfproviderdocs@latest
	go install github.com/katbyte/terrafmt@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.51.1

build: vet fmt
	go install

test: fmt
	go test -covermode atomic -coverprofile=covprofile -v $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4 -race

testacc: fmt
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -covermode atomic -coverprofile=covprofile 

lint:
	golangci-lint run ./...

vet:
	@echo "==> Running go vet ."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

