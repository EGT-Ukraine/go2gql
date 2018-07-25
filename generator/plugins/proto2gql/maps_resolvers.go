package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) mapResolverFunctionName(mapFile *parsedFile, mp *parser.Map) string {
	return "Resolve" + g.inputMapVariable(mapFile, mp)
}
func (g *Proto2GraphQL) fileInputMapResolvers(file *parsedFile) ([]graphql.MapInputObjectResolver, error) {
	var res []graphql.MapInputObjectResolver
	for _, msg := range file.File.Messages {
		for _, mapFld := range msg.MapFields {
			keyGoType, err := g.goTypeByParserType(mapFld.Map.KeyType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve key go type")
			}
			valueGoType, err := g.goTypeByParserType(mapFld.Map.ValueType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value go type")
			}
			keyTypeResolver, keyResolveWithErr, _, err := g.TypeValueResolver(file, mapFld.Map.KeyType, "")
			if err != nil {
				return nil, errors.Wrap(err, "failed to get key type value resolver")
			}
			valueParsedFile, err := g.parsedFile(mapFld.Map.ValueType.File())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message '%s' parsed file", dotedTypeName(msg.TypeName))
			}
			valueTypeResolver, valueResolveWithErr, _, err := g.TypeValueResolver(valueParsedFile, mapFld.Map.ValueType, "")
			if err != nil {
				return nil, errors.Wrap(err, "failed to get value type value resolver")
			}

			res = append(res, graphql.MapInputObjectResolver{
				FunctionName:           g.mapResolverFunctionName(file, mapFld.Map),
				KeyGoType:              keyGoType,
				ValueGoType:            valueGoType,
				KeyResolver:            keyTypeResolver,
				KeyResolverWithError:   keyResolveWithErr,
				ValueResolver:          valueTypeResolver,
				ValueResolverWithError: valueResolveWithErr,
			})
		}

	}
	return res, nil
}
