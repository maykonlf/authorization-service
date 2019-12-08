TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
GRPC_GATEWAY_PATH := $(shell go list -m -f '{{.Dir}}' all | grep 'github.com/grpc-ecosystem/grpc-gateway')
GRPC_GOOGLE_APIS_PATH := $(GRPC_GATEWAY_PATH)/third_party/googleapis/

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

certs:
	openssl genrsa -out server.key 2048
ifdef cn
	openssl req -new -x509 -sha256 -key server.key -out server.crt -subj "/CN=$(cn)"
else
	openssl req -new -x509 -sha256 -key server.key -out server.crt -subj "/CN=localhost"
endif

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

proto: $(PROTO_AUTHORIZATION_FILE_PATH) | $(PROTOC_GEN_GO) $(PROTOC)
	protoc -I api/proto/v1 -I "$(GRPC_GOOGLE_APIS_PATH)" -I "$(GRPC_GATEWAY_PATH)" --go_out=plugins=grpc:pkg/api/v1/ --grpc-gateway_out=logtostderr=true:pkg/api/v1/ --swagger_out=logtostderr=true:api/swagger/v1/ authorization.proto

.PHONY: default test cover fmt fmtcheck lint proto certs
