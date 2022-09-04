HERE=$(shell pwd)

.PHONY: default
default: build

.PHONY: build
.NOTPARALLEL: build
build: generate-protos dist/pb-create-archive.amd64 dist/pb-create-archive.arm64

.PHONY: clean
clean: 
	rm -rf dist
	make -C ../proto GO_OUT=${HERE}/proto clean

.PHONY: generate-protos
generate-protos:
	make -C ../proto GO_OUT=${HERE}/proto build-go

##########

dist/pb-create-archive.amd64: go.mod go.sum $(shell find . -name "*.go")
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -v -o dist/pb-create-archive.amd64 ./archive/cmd/pb-create-archive/.

dist/pb-create-archive.arm64: go.mod go.sum $(shell find . -name "*.go")
	mkdir -p dist
	GOOS=linux GOARCH=arm64 go build -v -o dist/pb-create-archive.arm64 ./archive/cmd/pb-create-archive/.
	