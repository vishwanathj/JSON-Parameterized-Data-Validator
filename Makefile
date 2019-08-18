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

all: deps build unit
build:
		$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(BINARY_SRC_DIR)/$(BINARY_NAME)/main.go
unit:
		#The idiomatic way to disable test caching explicitly is to use -count=1.
		TEST=true $(GOTEST) -v ./... -count=1 -tags=unit -coverprofile results/coverage_unit.out &> results/dbg_unit.out
		go tool cover -html=results/coverage_unit.out -o results/coverage_unit.html
		go tool cover -func=results/coverage_unit.out -o results/func_coverage.out
#integration:
#		$(GOTEST) -v ./... -tags=integration
display_unit_html:
		go tool cover -html=results/coverage_unit.out
display_int_html:
		go tool cover -html=results/coverage_integration.out
clean:
		$(GOCLEAN)
		rm -f $(GOPATH)/bin/$(BINARY_NAME)
		rm -f $(GOPATH)/bin/$(BINARY_UNIX)
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
		golint ./... &> results/lint.out
race:
		$(GOTEST) -race ${PKG_LIST}
msan:
		$(GOTEST) -msan -short ${PKG_LIST}


docker-build:
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml build
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.alpine.yml build
docker-publish-image:
		docker tag vnfdservice vishwanathj/vnfdservice
		docker push vishwanathj/vnfdservice
		docker tag vnfdservice_alpine vishwanathj/vnfdservice_alpine
		docker push vishwanathj/vnfdservice_alpine
docker-run:
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml up
docker-run-https:
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.https.yml up
docker-build-test:
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml build
docker-unit: docker-clean
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml up
docker-integration: docker-clean
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml -f $(DEPLOYMENT_DIR)/docker-compose.int.yml up
docker-ELK: docker-clean
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml -f $(DEPLOYMENT_DIR)/docker-compose.ELK.yml up -d
docker-jmeter: docker-clean
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml -f $(DEPLOYMENT_DIR)/docker-compose.jmeter.yml up
		#docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml -f $(DEPLOYMENT_DIR)/docker-compose.int.yml -f $(DEPLOYMENT_DIR)/docker-compose.jmeter.yml up

docker-clean:
		#docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml -f $(DEPLOYMENT_DIR)/docker-compose.int.yml -f $(DEPLOYMENT_DIR)/docker-compose.jmeter.yml down
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml -f $(DEPLOYMENT_DIR)/docker-compose.ELK.yml down
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml -f $(DEPLOYMENT_DIR)/docker-compose.jmeter.yml down
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml -f $(DEPLOYMENT_DIR)/docker-compose.int.yml down
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.unit.yml down
		docker-compose -f $(DEPLOYMENT_DIR)/docker-compose.yml down
		docker system prune -f
		docker volume prune -f