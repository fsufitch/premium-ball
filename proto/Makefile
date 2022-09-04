HERE=$(shell pwd)
GO_OUT:=${HERE}/dist/go
GO_OUT_FILES=$(shell F=$$(find ${GO_OUT} -name "*.go"); echo $${F:-${GO_OUT}/placeholder.go} )

PROTO_SOURCES=$(shell find . -name "*.proto")

.PHONY: all
all: build

.PHONY: build build-go build-ts
build: build-go build-ts
build-go: ${GO_OUT_FILES}
build-ts: build/ts/*

.PHONY: clean clean-go clean-ts
clean: clean-go clean-ts
clean-go:
	rm -rf ${GO_OUT_FILES}
clean-ts:
	true  # not implemented yet

##########

${GO_OUT_FILES}: ${PROTO_SOURCES}
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	mkdir -p ${GO_OUT}
	PATH=$$(go env GOPATH)/bin:$$PATH protoc -I=. --go_out=${GO_OUT} --go_opt=paths=source_relative ${PROTO_SOURCES}

build/ts/*: *.proto
	true  # not implemented yet
