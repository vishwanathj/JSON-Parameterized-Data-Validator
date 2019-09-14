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

all: deps unit
unit:
		mkdir -p $(TEST_RESULTS_DIR)
		#The idiomatic way to disable test caching explicitly is to use -count=1.
		$(GOTEST) -v ./... -count=1 -tags=unit -coverprofile $(TEST_RESULTS_DIR)/coverage_unit.out &> $(TEST_RESULTS_DIR)/dbg_unit.out
		go tool cover -html=$(TEST_RESULTS_DIR)/coverage_unit.out -o $(TEST_RESULTS_DIR)/coverage_unit.html
		go tool cover -func=$(TEST_RESULTS_DIR)/coverage_unit.out -o $(TEST_RESULTS_DIR)/func_coverage.out
display_unit_html:
		go tool cover -html=$(TEST_RESULTS_DIR)/coverage_unit.out
clean:
		docker system prune -f
		rm -rf $(TEST_RESULTS_DIR)
deps:
		dep ensure
		dep status
container_test:
		docker build -t vishwanathj/$(BINARY_NAME)_int -f $(BUILD_DIR)/Dockerfile_unit_test .
lint:
		golangci-lint --version; \
		golangci-lint run ./... --verbose
