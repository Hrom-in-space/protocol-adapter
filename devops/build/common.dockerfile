FROM golang:1.18-bullseye AS build
ENV PATH /usr/local/go/bin:$PATH


ENV GITHUB_PROXY="https://github.com/"

ENV PROTOC_VERSION="21.1"
ENV PROTOC_ARCH="x86_64"
ENV PROTOC_OS="linux"
ENV PROTOC_DOWNLOAD_URL="${GITHUB_PROXY}protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-${PROTOC_OS}-${PROTOC_ARCH}.zip"

ENV PROTOC_GEN_GO_VERSION="1.28.0"
ENV PROTOC_GEN_GO_ARCH="amd64"
ENV PROTOC_GEN_GO_OS="linux"
ENV PROTOC_GEN_GO_DOWNLOAD_URL="${GITHUB_PROXY}protocolbuffers/protobuf-go/releases/download/v${PROTOC_GEN_GO_VERSION}/protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.${PROTOC_GEN_GO_OS}.${PROTOC_GEN_GO_ARCH}.tar.gz"

ENV PROTOC_GEN_GO_GRPC_VERSION="1.2.0"
ENV PROTOC_GEN_GO_GRPC_ARCH="amd64"
ENV PROTOC_GEN_GO_GRPC_OS="linux"
ENV PROTOC_GEN_GO_GRPC_DOWNLOAD_URL="${GITHUB_PROXY}grpc/grpc-go/releases/download/cmd%2Fprotoc-gen-go-grpc%2Fv${PROTOC_GEN_GO_GRPC_VERSION}/protoc-gen-go-grpc.v${PROTOC_GEN_GO_GRPC_VERSION}.${PROTOC_GEN_GO_GRPC_OS}.${PROTOC_GEN_GO_GRPC_ARCH}.tar.gz"

ENV MINIMOCK_VERSION="3.0.10"
ENV MINIMOCK_ARCH="amd64"
ENV MINIMOCK_OS="linux"
ENV MINIMOCK_DOWNLOAD_URL="${GITHUB_PROXY}gojuno/minimock/releases/download/v${MINIMOCK_VERSION}/minimock_${MINIMOCK_VERSION}_${MINIMOCK_OS}_${MINIMOCK_ARCH}.tar.gz"

ENV GRPC_HEALTH_PROBE_VERSION="0.4.11"
ENV GRPC_HEALTH_PROBE_ARCH="amd64"
ENV GRPC_HEALTH_PROBE_OS="linux"
ENV GRPC_HEALTH_PROBE_DOWNLOAD_URL="${GITHUB_PROXY}grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-${GRPC_HEALTH_PROBE_OS}-${GRPC_HEALTH_PROBE_ARCH}"

RUN apt update && \
    apt -y install curl unzip git

WORKDIR /app

# protoc
RUN curl -L -o /tmp/protoc.zip ${PROTOC_DOWNLOAD_URL} && \
    unzip -j /tmp/protoc.zip bin/protoc -d /bin/
# protoc-gen-go
RUN curl -L -o /tmp/protoc-gen-go.tar.gz ${PROTOC_GEN_GO_DOWNLOAD_URL} && \
    tar -xvf /tmp/protoc-gen-go.tar.gz -C /bin/
# protoc-gen-go-grpc
RUN curl -L -o /tmp/protoc-gen-go-grpc.tar.gz ${PROTOC_GEN_GO_GRPC_DOWNLOAD_URL} && \
    tar -C /bin/ -xvf /tmp/protoc-gen-go-grpc.tar.gz ./protoc-gen-go-grpc
# minimock
RUN curl -L -o /tmp/minimock.tar.gz ${MINIMOCK_DOWNLOAD_URL} && \
    tar -C /bin/ -xvf /tmp/minimock.tar.gz minimock

RUN curl -L -o /app/grpc_health_probe ${GRPC_HEALTH_PROBE_DOWNLOAD_URL} && \
    chmod +x /app/grpc_health_probe

COPY go.mod go.sum main.go Makefile ./
COPY  internal/ ./internal
COPY  proto/ ./proto
COPY  vendor/ ./vendor

RUN make gen
RUN make build


FROM debian:bullseye

ENV GOMAXPROCS=1

RUN apt update && apt install -y wget ca-certificates

WORKDIR /app
COPY --from=build /app/protocol-adapter /app/protocol-adapter
COPY --from=build /app/grpc_health_probe /app/grpc_health_probe

RUN groupadd -r app && useradd -r -m -g app app && chown -R app:app /app
USER app

CMD ["./protocol-adapter"]
