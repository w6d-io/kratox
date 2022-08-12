SHELL=/bin/bash -o pipefail

export GO111MODULE        := on
export PATH               := bin:${PATH}
export PWD                := $(shell pwd)
export BUILD_DATE         := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
export VCS_REF            := $(shell git rev-parse HEAD)
export NEXT_TAG           ?=

ifeq (,$(shell go env GOOS))
GOOS       = $(shell echo $OS)
else
GOOS       = $(shell go env GOOS)
endif


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

.PHONY: bin/goreadme
bin/goreadme:
	GOBIN=${CURDIR}/bin \
	go install github.com/posener/goreadme/cmd/goreadme

.PHONY: readme
readme: bin/goreadme
	./third_party/script/create_readme.sh


.PHONY: kratos
KRATOS_BINARY = $(shell pwd)/bin/kratos
KRATOS_TAR.GZ = $(shell pwd)/bin/kratos.tar.gz
SCRIPTBASH = $(shell pwd)/makefile.sh
GOBIN = $(shell pwd)/bin
ifeq (darwin,$(GOOS))
KRATOS_BINARY_URL=https://github.com/ory/kratos/releases/download/v0.10.1/kratos_0.10.1-macOS_sqlite_64bit.tar.gz
else
KRATOS_BINARY_URL=https://github.com/ory/kratos/releases/download/v0.10.1/kratos_0.10.1-linux_sqlite_64bit.tar.gz
endif
kratos: start ##init kratos
ifeq (,$(wildcard $(KRATOS_BINARY)))
	mkdir -p $(GOBIN)
	wget $(KRATOS_BINARY_URL) -O $(KRATOS_TAR.GZ)
	tar -xf $(KRATOS_TAR.GZ) -C $(GOBIN)
	chmod +x $(KRATOS_BINARY)
	mkdir -p /var/lib/sqlite
	$(SCRIPTBASH) config &
else
	$(info ************ BINARY ALREADY EXIST **********)
endif

start: ## if binary file cannot execute verify the KRATOS_BINARY_URL OS (default linux or macosx)
	nohup $(KRATOS_BINARY) serve --dev -c $(GOBIN)/kratos.yml &

stop:
ifeq (darwin,$(GOOS))
	lsof -i -P | grep 4434 | sed -e 's/.*kratos     *//' -e 's#/.*##' | sed 's/ .*//' | xargs kill
else
	netstat -lnp | grep 4434 | sed -e 's/.*LISTEN *//' -e 's#/.*##' | xargs kill
endif

clean:
	rm -rf bin