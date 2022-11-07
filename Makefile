proto:  ## Generate client/server from proto
	protoc \
	--proto_path=proto \
	--go_out=./internal/generated/common \
	--go-grpc_out=./internal/generated/common \
	error.proto && \
	for flname in $(shell ls ./proto/protocol-adapter); do \
		protoc \
		--proto_path=proto \
		--go_out=./internal/generated/server \
		--go-grpc_out=./internal/generated/server \
		--go_opt=Merror.proto=protocol-adapter/internal/generated/common/error \
		protocol-adapter/$$flname; \
	done

gen:  ##
	go generate ./...

run:
	go run ./cmd/protocol-adapter

test:
	go test ./...


coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	golangci-lint run

lint-docker:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run

lint-fix:
	golangci-lint run --fix

build:
	go build \
	-ldflags "-X protocol-adapter/internal/config.BuildVersion=${VERSION}" \
	-o protocol-adapter-${VERSION}-${GOOS}-${GOARCH} \
	./cmd/protocol-adapter

mock:
	go run cmd/mock-server/main.go ./fixtures
