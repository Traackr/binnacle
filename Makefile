.PHONY: build clean docker-build docker-test-coverage docker-test-unit .docker-prep tag-version test-coverage test-unit tools version version-bump

BUILD_IMAGE=golang:1.8.1
COVERAGE_DIR=./$(DOCS_DIR)/coverage
DOCS_DIR=./docs
EXTERNAL_TOOLS=\
	github.com/kardianos/govendor \
	github.com/mitchellh/gox \
	golang.org/x/tools/cmd/cover \
	github.com/axw/gocov/gocov \
	gopkg.in/matm/v1/gocov-html \
	github.com/golang/lint/golint \
	github.com/tebeka/go2xunit \
	github.com/ugorji/go/codec/codecgen

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LOCAL_TARGET ?= "$(GOOS)_$(GOARCH)"
PACKAGE := $(shell pwd | awk -F/ '{print $$NF}')
PACKAGE_PATH := $(subst ${GOPATH}/src/,,$(shell pwd))
UNIT_DIR=./$(DOCS_DIR)/unit
VERSION_FILE := "VERSION"
VERSION := $(shell cat $(VERSION_FILE))

build:
	@scripts/build.sh

clean:
	@rm -rf bin pkg

docker-build: .docker-prep
	@echo "==> Starting docker container for building..."
	@docker run --rm \
		-v "$$PWD":/go/src/$(PACKAGE_PATH) \
		-w /go/src/$(PACKAGE_PATH) \
		-e GENERATE_PACKAGES \
		-e LOCAL_TARGET=$(LOCAL_TARGET) \
		-e TARGETS \
		$(BUILD_IMAGE) \
		make tools build
	
.docker-prep:
	@echo "==> Pulling $(BUILD_IMAGE)..."
	@docker pull $(BUILD_IMAGE) > /dev/null

docker-test-coverage: .docker-prep
	@echo "==> Starting docker container for testing..."
	@docker run --rm \
		-v "$$PWD":/go/src/$(PACKAGE_PATH) \
		-w /go/src/$(PACKAGE_PATH) \
		$(BUILD_IMAGE) \
		make tools test-coverage

docker-test-unit: .docker-prep
	@echo "==> Starting docker container for testing..."
	@docker run --rm \
		-v "$$PWD":/go/src/$(PACKAGE_PATH) \
		-w /go/src/$(PACKAGE_PATH) \
		$(BUILD_IMAGE) \
		make tools test-unit

test-coverage:
	@mkdir -p $(COVERAGE_DIR)
	@rm -rf $(COVERAGE_DIR)/*
	@touch $(COVERAGE_DIR)/coverage.tmp
	@echo 'mode: atomic' > $(COVERAGE_DIR)/coverage.txt
	@go list ./... | grep -v vendor/ | grep -v mock | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage.tmp {} && tail -n +2 $(COVERAGE_DIR)/coverage.tmp >> $(COVERAGE_DIR)/coverage.txt'
	@rm $(COVERAGE_DIR)/coverage.tmp
	@go tool cover -html=$(COVERAGE_DIR)/coverage.txt -o $(COVERAGE_DIR)/report.html
	@gocover-cobertura < $(COVERAGE_DIR)/coverage.txt > $(COVERAGE_DIR)/coverage.xml
	@rm $(COVERAGE_DIR)/coverage.txt

test-unit:
	@mkdir -p $(UNIT_DIR)
	@go test -v $$(go list ./... | grep -v vendor/ | grep -v mock) | tee $(UNIT_DIR)/test.out
	@go2xunit -fail -input $(UNIT_DIR)/test.out -output $(UNIT_DIR)/xunit.xml
	@rm $(UNIT_DIR)/test.out

tools:
	@for tool in $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
		go get $$tool; \
	done

travis: test-unit build

version:
	@if [ -e $(VERSION_FILE) ]; then \
		echo "$$VERSION"; \
	else \
		echo "No version file found"; \
		exit 1; \
	fi;

version-bump:
	@if [ -e $(VERSION_FILE) ]; then \
		MAJOR=`echo "$$VERSION" | cut -d. -f 1`; \
		MINOR=`echo "$$VERSION" | cut -d. -f 2`; \
		PATCH=`echo "$$VERSION" | cut -d. -f 3`; \
		NEW_PATCH=`echo "$$(( $$PATCH + 1 ))"`; \
		echo "$$MAJOR.$$MINOR.$$NEW_PATCH" > $(VERSION_FILE); \
	else \
		echo "No version file found"; \
		exit 1; \
	fi;

version-push:
	@if [ -e $(VERSION_FILE) ]; then \
		git commit -a -m "[$(PACKAGE)] Version Bump"; \
		git pull --rebase origin master; \
		git push origin HEAD:master; \
	else \
		echo "No version file found"; \
		exit 1; \
	fi;

version-tag:
	@if [ -e $(VERSION_FILE) ]; then \
		git tag -f -a -m "[$(PACKAGE)] Version $(VERSION)" "$$VERSION"; \
		git pull --rebase origin master; \
		git push origin HEAD:master --force; \
	else \
		echo "No version file found"; \
		exit 1; \
	fi;

vet:
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
