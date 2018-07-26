default: build_templates install

install:
	@go install ./cmd/go2gql/

build:
	@go build -o ./bin/go2gql ./cmd/go2gql

build_templates:
	go-bindata -prefix ./generator/plugins/graphql -o ./generator/plugins/graphql/templates.go -pkg graphql ./generator/plugins/graphql/templates
	go-bindata -prefix ./generator/plugins/swagger2gql -o ./generator/plugins/swagger2gql/templates.go -pkg swagger2gql ./generator/plugins/swagger2gql/templates

test:
	go test ./...


.PHONY: install


