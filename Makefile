# Go parameters
# Reference --> https://sohlich.github.io/post/go_makefile/
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=json-data-validator
BINARY_UNIX=$(BINARY_NAME)
BINARY_SRC_DIR=cmd
BUILD_DIR=build/package
DEPLOYMENT_DIR=deployments/docker-compose
#https://medium.com/pantomath/go-tools-gitlab-how-to-do-continuous-integration-like-a-boss-941a3a9ad0b6
PKG_LIST=$(shell go list ./... | grep -v /vendor/)
TEST_RESULTS_DIR=test_results

all: deps build unit
build:
		$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(BINARY_SRC_DIR)/$(BINARY_NAME)/main.go
unit:
		#The idiomatic way to disable test caching explicitly is to use -count=1.
		mkdir -p $(TEST_RESULTS_DIR)
		$(GOTEST) -v ./... -count=1 -tags=unit -coverprofile $(TEST_RESULTS_DIR)/coverage_unit.out &> $(TEST_RESULTS_DIR)/dbg_unit.out
		go tool cover -html=$(TEST_RESULTS_DIR)/coverage_unit.out -o $(TEST_RESULTS_DIR)/coverage_unit.html
		go tool cover -func=$(TEST_RESULTS_DIR)/coverage_unit.out -o $(TEST_RESULTS_DIR)/func_coverage.out
display_unit_html:
		go tool cover -html=$(TEST_RESULTS_DIR)/coverage_unit.out
clean:
		rm -rf $(TEST_RESULTS_DIR)
		#$(GOCLEAN)
		#rm -f $(GOPATH)/bin/$(BINARY_NAME)
		#rm -f $(GOPATH)/bin/$(BINARY_UNIX)
run:
		$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(BINARY_SRC_DIR)/$(BINARY_NAME)/main.go
		$(BINARY_NAME)
deps:
		$(GOGET) gopkg.in/tomb.v2 github.com/stretchr/testify golang.org/x/lint/golint github.com/t-yuki/gocover-cobertura
		$(GOGET) -d -v ./...
container:
		docker build -t vishwanathj/$(BINARY_NAME) -f $(BUILD_DIR)/Dockerfile .
container_test:
		docker build -t vishwanathj/$(BINARY_NAME)_int -f $(BUILD_DIR)/Dockerfile_test .
lint:
		#golint ./... &> $(TEST_RESULTS_DIR)/lint.out
		golangci-lint --version; \
		golangci-lint run ./... --verbose
race:
		$(GOTEST) -race ${PKG_LIST}
msan:
		$(GOTEST) -msan -short ${PKG_LIST}
