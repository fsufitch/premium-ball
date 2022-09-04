####################
### Builder image that creates the backend binary artifacts 
FROM ubuntu:22.04 AS backend-binaries-builder

RUN apt-get update && apt-get install -y curl make protobuf-compiler

ARG GO_DIST=https://go.dev/dl/go1.19.linux-amd64.tar.gz
RUN mkdir -p /opt && curl -L "${GO_DIST}" > /opt/go.tar.gz
RUN tar xf /opt/go.tar.gz -C /opt
ENV PATH=/opt/go/bin:$PATH

WORKDIR /pball
COPY backend backend
COPY proto proto

RUN make -C backend


####################
### Scratch image containing *only* the backend binary artifacts
FROM scratch AS backend-binaries

COPY --from=backend-binaries-builder \
    /pball/backend/dist/pb-create-archive.amd64 \
    /pball/backend/dist/pb-create-archive.arm64 \
    /


####################
### Builder image for downloading/archiving the pokemontcg.io card database
FROM ubuntu:22.04 AS card-archive-builder

WORKDIR /pball
COPY --from=docker.io/fsufitch/pokemon-premium-ball-backend-binaries \
    /pb-create-archive.amd64 \
    /pb-create-archive.arm64 \
    .

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates && \
    if    [ "$(uname -m)" = "x86_64" ]; then ./pb-create-archive.amd64 cards.zip; \
    elif  [ "$(uname -m)" = "arm64" ]; then ./pb-create-archive.arm64 cards.zip; \
    else echo "ERROR: host is not amd64 or arm64" >/dev/stderr; exit 1; \
    fi


####################
### Scratch image containing *only* the card archive
FROM scratch AS card-archive
COPY --from=card-archive-builder /pball/cards.zip /


####################
### Dummy to have a fast default; use --target to target stuff
FROM scratch AS default