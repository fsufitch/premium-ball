VERSION:=0.1
PULL_VERSION:=latest

PODMAN:=$(shell which podman || which docker || echo docker-podman-not-found)
TAG_PREFIX:=docker.io/fsufitch/premium-ball

.DEFAULT: help
.PHONY: help
help:
	@echo Available targets:
	@echo - pull
	@echo - backend-binaries 
	@echo - push-backend-binaries
	@echo - card-archive
	@echo - push-card-archive
	@echo
	@echo Configuration:
	@echo - VERSION=${VERSION}
	@echo - PULL_VERSION=${PULL_VERSION}
	@echo - PODMAN=${PODMAN}
	@echo - TAG_PREFIX=${TAG_PREFIX}

.PHONY: pull
pull:
	${PODMAN} pull \
		${TAG_PREFIX}.backend-binaries:${PULL_VERSION} \
		${TAG_PREFIX}.card-archive:${PULL_VERSION} \


.PHONY: backend-binaries
backend-binaries:
	${PODMAN} build \
		-f Containerfile \
		-t ${TAG_PREFIX}.backend-binaries:${VERSION} \
		--target backend-binaries \
		.

.PHONY: push-backend-binaries
push-backend-binaries:
	${PODMAN} tag ${TAG_PREFIX}.backend-binaries:${VERSION}  ${TAG_PREFIX}.backend-binaries:latest
	${PODMAN} push ${TAG_PREFIX}.backend-binaries:${VERSION}
	${PODMAN} push ${TAG_PREFIX}.backend-binaries:latest

.PHONY: card-archive
card-archive:
	${PODMAN} build \
		-f Containerfile \
		-t ${TAG_PREFIX}.card-archive:${VERSION} \
		--target card-archive \
		.

.PHONY: push-card-archive
push-card-archive:
	${PODMAN} tag ${TAG_PREFIX}.card-archive:${VERSION} ${TAG_PREFIX}.card-archive:latest
	${PODMAN} push ${TAG_PREFIX}.card-archive:${VERSION}
	${PODMAN} push ${TAG_PREFIX}.card-archive:latest
