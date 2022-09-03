.PHONY: default
default: all

.PHONY: all
all: generate-protos build

.PHONY: build
build: build/pb-create-archive

.PHONY: clean
clean: 
	rm -rf build
	make -C ../proto GO_OUT=${PWD}/proto clean

.PHONY: generate-protos
generate-protos:
	make -C ../proto GO_OUT=${PWD}/proto build-go

##########

build/pb-create-archive: go.mod go.sum $(shell find . -name "*.go")
	mkdir -p build
	go build -o build/pb-create-archive ./archive/cmd/pb-create-archive/.