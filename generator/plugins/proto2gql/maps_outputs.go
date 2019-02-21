package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) outputMapGraphQLName(mapFile *parsedFile, res *parser.Map) string {
	return g.outputMessageVariable(mapFile, res.Message) + "__" + res.Field.Name
}

func (g *Proto2GraphQL) outputMapVariable(mapFile *parsedFile, res *parser.Map) string {
	return g.outputMessageVariable(mapFile, res.Message) + "__" + res.Field.Name
}

func (g *Proto2GraphQL) fileMapOutputObjects(file *parsedFile) ([]graphql.MapOutputObject, error) {
	var res []graphql.MapOutputObject
	for _, msg := range file.File.Messages {
		for _, mapFld := range msg.MapFields {
			keyTypResolver, err := g.TypeOutputGraphQLTypeResolver(file, mapFld.Map.KeyType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve key input type resolver")
			}
			valueFile, err := g.parsedFile(mapFld.Map.ValueType.File())
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value type file")
			}

			valueResolver := func(arg string, ctx graphql.BodyContext) string {
				return `src := p.Source.(map[string]interface{})
				if src == nil {
					return nil, nil
				}
				return src["value"], nil`
			}

			valueTypResolver, err := g.TypeOutputGraphQLTypeResolver(valueFile, mapFld.Map.ValueType)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value input type resolver")
			}

			valueMessage, ok := mapFld.Map.ValueType.(*parser.Message)

			if ok {
				msgCfg, err := valueFile.Config.MessageConfig(valueMessage.Name)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve message %s config", valueMessage.Name)
				}

				if msgCfg.UnwrapField {
					fields := valueMessage.GetFields()
					if len(fields) != 1 {
						return nil, errors.Errorf("can't unwrap %s output message because it contains more that 1 field", valueMessage.Name)
					}

					unwrappedField := fields[0]

					responseGoType, err := g.goTypeByParserType(valueMessage)

					if err != nil {
						return nil, err
					}

					valueResolver = func(arg string, ctx graphql.BodyContext) string {
						return `src := p.Source.(map[string]interface{})
						if src == nil {
							return nil, nil
						}

						return src["value"].(` + responseGoType.String(ctx.Importer) + `).` + camelCase(unwrappedField.GetName()) + `, nil`
					}

					valueTypResolver, err = g.TypeOutputGraphQLTypeResolver(valueFile, unwrappedField.GetType())
					if err != nil {
						return nil, errors.Wrap(err, "failed to resolve value input type resolver")
					}
					if unwrappedField.IsRepeated() {
						valueTypResolver = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(valueTypResolver))
					}
				}
			}

			res = append(res, graphql.MapOutputObject{
				VariableName:    g.outputMapVariable(file, mapFld.Map),
				GraphQLName:     g.outputMapGraphQLName(file, mapFld.Map),
				KeyObjectType:   keyTypResolver,
				ValueObjectType: valueTypResolver,
				ValueResolver:   valueResolver,
			})
		}
	}
	return res, nil
}
