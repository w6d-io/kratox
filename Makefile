SHELL=/bin/bash -o pipefail

export GO111MODULE        := on
export PATH               := bin:${PATH}
export PWD                := $(shell pwd)
export BUILD_DATE         := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export VCS_REF            := $(shell git rev-parse HEAD)
export NEXT_TAG           ?=

GO_DEPENDENCIES = golang.org/x/tools/cmd/goimports

define make-go-dependency
  # go install is responsible for not re-building when the code hasn't changed
  bin/$(notdir $1): go.mod go.sum Makefile
	GOBIN=$(PWD)/bin/ go install $1
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency, $(dep))))
$(call make-lint-dependency)

# Formats the code
.PHONY: format
format:
	goimports -w -local github.com/w6d-io,gitlab.w6d.io/w6d .

.PHONY: changelog
changelog:
	git-chglog -o CHANGELOG.md --next-tag $(NEXT_TAG)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

