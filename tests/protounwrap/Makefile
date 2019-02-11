dataloader_tests:
	rm -rf generated/*
	mkdir -p generated/clients
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/items.proto
	mockgen -destination=mock/item.go -package=mock github.com/EGT-Ukraine/go2gql/tests/protounwrap/generated/clients/apis ItemsServiceClient
	go run ../../cmd/go2gql/main.go ../../cmd/go2gql/basic_plugins.go
