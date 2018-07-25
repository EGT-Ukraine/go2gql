package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) inputMapGraphQLName(mapFile *parsedFile, res *parser.Map) string {
	return g.inputMessageVariable(mapFile, res.Message) + "__" + camelCase(res.Field.Name)
}

func (g *Proto2GraphQL) inputMapVariable(mapFile *parsedFile, res *parser.Map) string {
	return g.inputMessageVariable(mapFile, res.Message) + "__" + camelCase(res.Field.Name)
}

func (g *Proto2GraphQL) fileMapInputObjects(file *parsedFile) ([]graphql.MapInputObject, error) {
	var res []graphql.MapInputObject
	for _, msg := range file.File.Messages {
		for _, mapFld := range msg.MapFields {
			keyTypResolver, err := g.TypeInputTypeResolver(file, mapFld.Map.KeyType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve key input type resolver")
			}
			valueFile, err := g.parsedFile(mapFld.Map.ValueType.File())
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value type file")
			}
			valueTypResolver, err := g.TypeInputTypeResolver(valueFile, mapFld.Map.ValueType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value input type resolver")
			}

			res = append(res, graphql.MapInputObject{
				VariableName:    g.inputMapVariable(file, mapFld.Map),
				GraphQLName:     g.inputMapGraphQLName(file, mapFld.Map),
				KeyObjectType:   keyTypResolver,
				ValueObjectType: valueTypResolver,
			})
		}

	}
	return res, nil
}
