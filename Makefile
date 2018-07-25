default: build_templates install

install:
	@go install ./cmd/proto2gql/

build:
	@go build -o ./bin/proto2gql ./cmd/proto2gql

build_templates:
	go-bindata -prefix ./generator/plugins/graphql -o ./generator/plugins/graphql/templates.go -pkg graphql ./generator/plugins/graphql/templates
	go-bindata -prefix ./generator/plugins/swagger2gql -o ./generator/plugins/swagger2gql/templates.go -pkg swagger2gql ./generator/plugins/swagger2gql/templates

test:
	go test ./...


.PHONY: install


