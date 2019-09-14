# Go parameters
# Reference --> https://sohlich.github.io/post/go_makefile/
GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=json-data-validator
BUILD_DIR=build/package
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
