default: proto

proto:
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:${GOPATH}/src  test.proto
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:${GOPATH}/src  common/common.proto
	@protoc -I=${GOPATH}/src:. --go_out=plugins=grpc:${GOPATH}/src  test_scope.proto
	@go run ../cmd/go2gql/main.go ../cmd/go2gql/basic_plugins.go

.PHONY: proto
