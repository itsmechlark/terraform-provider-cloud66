TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=cloud66
VERSION=$(shell git describe --tags --always)

default: build

build: vet fmt
	go install

test: fmt
	go test -covermode atomic -coverprofile=covprofile -v $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4 -race

testacc: fmt
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -covermode atomic -coverprofile=covprofile 

lint: tools terraform-provider-lint golangci-lint

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

