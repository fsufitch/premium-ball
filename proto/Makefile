HERE=$(shell pwd)
GO_OUT:=${HERE}/build/go
GO_OUT_FILES:=$(shell F=$$(find ${GO_OUT} -name "*.go"); echo $${F:-${GO_OUT}/placeholder.go} )

.PHONY: all
all: build

.PHONY: build build-go build-ts
build: build-go build-ts
build-go: build/pball-protoc.timestamp ${GO_OUT_FILES}
build-ts: build/pball-protoc.timestamp build/ts/*

.PHONY: clean
clean:
	rm -rf ${GO_OUT_FILES} build/ts

.PHONY: clean-all
clean-all: clean
	podman rmi -f pball-protoc


##########

build/pball-protoc.timestamp: Containerfile.pball-protoc
	podman build -f Containerfile.pball-protoc -t pball-protoc
	mkdir -p build && touch build/pball-protoc.timestamp

${GO_OUT_FILES}: $(shell find . -name "*.proto")
	mkdir -p ${GO_OUT}
	podman run \
		-v ${HERE}:/proto_src \
		-v ${GO_OUT}:/proto_out \
		pball-protoc sh -c 'protoc \
			-I=/proto_src \
			--go_out=/proto_out \
			--go_opt=paths=source_relative \
			$$(find /proto_src/ -name "*.proto")'


build/ts/*: *.proto
	echo ts build not implemented