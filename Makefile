# Go parameters
# Reference --> https://sohlich.github.io/post/go_makefile/
GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=json-data-validator
BUILD_DIR=build/package
TEST_RESULTS_DIR=test_results
LINT_DKR_IMG=golangci/golangci-lint:v1.18.0
GOSEC_VER=v2.2.0
#LINT_DKR_IMG=golangci/golangci-lint:v1.18.0

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
lint_dkr:
		docker run --rm -v ${PWD}:/go/src/github.com/vishwanathj/JSON-Parameterized-Data-Validator -w /go/src/github.com/vishwanathj/JSON-Parameterized-Data-Validator $(LINT_DKR_IMG) \
		sh -c "go get -u github.com/golang/dep/cmd/dep  && dep ensure -v && golangci-lint run -v"
docker-gosec:
		#docker run -it -v <YOUR PROJECT PATH>/<PROJECT>:/<PROJECT> securego/gosec /<PROJECT>/...
		docker run -it -v ${PWD}:/JSONPDV securego/gosec /JSONPDV/...
gosec:
		curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH) v2.2.0
		$(GOPATH)/gosec ./...