# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
GOOS ?= linux
GOARCH ?= amd64
GO     ?= GO15VENDOREXPERIMENT=1 go
GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
pkgs         = $(shell $(GO) list ./... | grep -v /vendor/)
PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME  ?= mercury200hg/metrics-server-prometheus-exporter
DOCKER_IMAGE_TAG ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKERFILE ?= Dockerfile
all: vendor format build docker
vendor:
	@echo ">> Adding vendors"
	@$(GO) mod vendor
format: 
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)
vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)
build:
	@echo ">> Building project metrics-server-prometheus-exporter"
	@$(GO) get
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build
docker:
	@echo ">> Building Docker image $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)"
	@docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	@docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)