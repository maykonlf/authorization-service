TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
GRPC_GATEWAY_LIBS := $(shell go list -m -f '{{.Dir}}' all | grep grpc-gateway)/third_party/googleapis/
AUTHORIZATION_PROTO_FILE_PATH := proto/authorization/v1/authorization.proto

default: test

fmt:
	gofmt -w $(GOFMT_FILES)

lint:
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --tests=false --skip-files=mock.go --disable=goimports --enable-all

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

cover: fmtcheck lint
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

authorization.pb.go: $(PROTO_AUTHORIZATION_FILE_PATH) | $(PROTOC_GEN_GO) $(PROTOC)
	protoc -I . -I $(GRPC_GATEWAY_LIBS) --grpc-gateway_out=logtostderr=true:. --go_out=plugins=grpc:. $(AUTHORIZATION_PROTO_FILE_PATH)

proto: authorization.pb.go

.PHONY: default test cover fmt fmtcheck lint
