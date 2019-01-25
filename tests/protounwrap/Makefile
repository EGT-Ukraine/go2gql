dataloader_tests:
	rm -rf generated/*
	mkdir -p generated/clients
	protoc --go_out=paths=source_relative,plugins=grpc:generated/clients apis/items.proto

	go run ../../cmd/go2gql/main.go ../../cmd/go2gql/basic_plugins.go
