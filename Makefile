# Go parameters
# Reference --> https://sohlich.github.io/post/go_makefile/
GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=json-data-validator
BUILD_DIR=build/package
# change to value of TEST_RESULTS_DIR would need a corresponding change in .circleci/config.yml file
TEST_RESULTS_DIR=$(HOME)/test_results
LINT_DKR_IMG=golangci/golangci-lint:v1.40.1
#LINT_DKR_IMG=golangci/golangci-lint:v1.23-alpine
GOSEC_VER=v2.2.0

all: deps unit
unit:
		mkdir -p $(TEST_RESULTS_DIR)
		#The idiomatic way to disable test caching explicitly is to use -count=1.
		#$(GOTEST) -v ./... -count=1 -tags=unit -coverprofile $(TEST_RESULTS_DIR)/coverage_unit.out &> $(TEST_RESULTS_DIR)/dbg_unit.out
		$(GOTEST) -v ./... -count=1 -tags=unit -coverprofile $(TEST_RESULTS_DIR)/coverage_unit.out
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
docker-unit-tests:
		docker run --rm -v ${PWD}:/go/src/github.com/JSONPDV -w /go/src/github.com/JSONPDV golang:1.15-buster go test -v ./... -count=1 -tags=unit
lint:
		golangci-lint --version; \
		golangci-lint run ./... --verbose
docker-lint:
		docker run --rm -v ${PWD}:/app -w /app $(LINT_DKR_IMG) golangci-lint run -v
docker-gosec:
		docker run --rm -it -v ${PWD}:/app securego/gosec:$(GOSEC_VER) /app/...
gosec:
		curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ~/tmp v2.2.0
		~/tmp/gosec ./...