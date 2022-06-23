# INTEL CONFIDENTIAL
#
# Copyright 2021-2021 Intel Corporation.
#
# This software and the related documents are Intel copyrighted materials, and your use of
# them is governed by the express license under which they were provided to you ("License").
# Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
# disclose or transmit this software or the related documents without Intel's prior written permission.
#
# This software and the related documents are provided as is, with no express or implied warranties,
# other than those that are expressly stated in the License.

export GO111MODULE = on
export GOPRIVATE="github.com/smart-edge-open/*"
export ACK_GINKGO_DEPRECATIONS=1.16.4

.PHONY: \
	clean \
	build \
	edgednssvr \
	edgednscli \
	lint test help
COPY_DOCKERFILES := $(shell /usr/bin/cp -rfT ./build/ ./dist/)
VER ?= 1.0

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo ""
	@echo "Helper targets:"
	@echo "  clean                  to clean build artifacts"
	@echo "  build                  to build both edgedns executables without images"
	@echo "  lint                   to run linter on Go code"
	@echo "  test                   to run tests on Go code"
	@echo "  test-cov               to run coverage tests on Go code"
	@echo "  help                   to show this message"
	@echo ""
	@echo "Single targets:"
	@echo "  edgednssvr             to build only docker image of the Edge DNS Service"
	@echo "  edgednscli      		to build edgednscli to the ./dist/ folder"
	
clean:
	rm -rf ./dist

test:
	http_proxy= https_proxy= HTTP_PROXY= HTTPS_PROXY= ginkgo -v -r  -gcflags=-l --randomizeSuites --failOnPending --skipPackage=vendor,k8s

test-cov:
	rm -rf coverage.out*
	http_proxy= https_proxy= HTTP_PROXY= HTTPS_PROXY= ginkgo -v -r --randomizeSuites --failOnPending --skipPackage=vendor,k8s \
	 -gcflags=-l -cover -coverprofile=coverage.out -outputdir=.
	sed '1!{/^mode/d;}' coverage.out > coverage.out.fix
	go tool cover -html=coverage.out.fix

edgednssvr:
	GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o ./dist/$@ ./cmd/$@
	docker build ./dist/ --build-arg username=edgednssvr -f ./build/Dockerfile -t $@:${VER}

edgednscli:
	go build -o dist/edgednscli ./cmd/edgednscli

build:
	GOPRIVATE="github.com/smart-edge-open/*" go get github.com/smart-edge-open/edge-services/common/log
	GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o ./dist/edgednssvr ./cmd/edgednssvr
	go build -o dist/edgednscli ./cmd/edgednscli

lint:
	golangci-lint run
