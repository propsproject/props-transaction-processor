OUT :=./build/bin/pending-props-tp
PKG := github.com/propsproject/props-transaction-processor/cmd
DOCKERFILE := ./build/package/Dockerfile
VERSION := $(shell git describe --tags --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: run

build-deploy:
	docker build -f ${DOCKERFILE} -t $(REPO):$(BUILD_NUMBER) .

docker-image:
	docker build -f ${DOCKERFILE} -t propsprojectservices/props-transaction-processor:latest .

deps:
	go get ${PKG}

build:
	go build -i -v -o ${OUT} -ldflags="-X main.repoVersion=${VERSION}" ${PKG}

test:
	go test ./...

vet:
	go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

runtp:
	./${OUT}-v${VERSION} --validator-endpoint $(validator)

out:
	@echo ${OUT}-v${VERSION}

protos:

	protoc -I ./protos ./protos/transaction.proto ./protos/events.proto ./protos/payload.proto ./protos/balance.proto ./protos/users.proto ./protos/activity.proto --go_out=./core/proto/pending_props_pb
	protoc -I ./protos ./protos/transaction.proto ./protos/events.proto ./protos/payload.proto ./protos/balance.proto ./protos/users.proto ./protos/activity.proto --js_out=import_style=commonjs,binary:./dev-cli/proto

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run protos runtp build docker-image vet lint out deps build-cgo