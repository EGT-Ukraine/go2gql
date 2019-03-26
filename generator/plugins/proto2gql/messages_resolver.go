package proto2gql

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) inputMessageResolverName(msgFile *parsedFile, message *parser.Message) string {
	return "Resolve" + g.inputMessageGraphQLName(msgFile, message)
}

func (g *Proto2GraphQL) oneOfValueAssigningWrapper(file *parsedFile, msg *parser.Message, field *parser.NormalField) graphql.AssigningWrapper {
	return func(arg string, ctx graphql.BodyContext) string {
		return "&" + ctx.Importer.Prefix(file.GRPCSourcesPkg) + camelCaseSlice(msg.TypeName) + "_" + camelCase(field.Name) + "{" + arg + "}"
	}
}

func (g *Proto2GraphQL) inputUnwrappedMessagesResolver(file *parsedFile, msg *parser.Message) (*graphql.InputObjectResolver, error) {
	msgGoType, err := g.goTypeByParserType(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve message go type")
	}

	if len(msg.NormalFields) != 1 {
		return nil, errors.Wrapf(err, "can't unwrap message %s. Must have 1 field", msg.Name)
	}

	fld := msg.NormalFields[0]

	goType, err := g.goTypeByParserType(fld.Type)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get go type by parser type")
	}

	return &graphql.InputObjectResolver{
		FunctionName: g.inputMessageResolverName(file, msg),
		OutputGoType: msgGoType,
		Fields: []graphql.InputObjectResolverField{
			{
				GraphQLInputFieldName: fld.Name,
				OutputFieldName:       camelCase(fld.Name),
				ValueResolver: func(arg string, ctx graphql.BodyContext) string {
					goTypeString := goType.String(ctx.Importer)

					// We can't use default graphql input objects resolver. So goType never be a slice.
					if fld.Repeated {
						return `func() []` + goTypeString + ` {
								in := i.([]interface{})
								res := make([]` + goTypeString + `, len(in))
								for i, val := range in {
									res[i] = ` + "val.(" + goTypeString + ")" + `
								}

								return res
							}()
						`
					}

					if _, ok := fld.Type.(*parser.Enum); ok {
						return goTypeString + "(i.(int))"
					}

					return "i.(" + goTypeString + ")"
				},
				GoType: goType,
			},
		},
	}, nil
}

func (g *Proto2GraphQL) fileInputMessagesResolvers(file *parsedFile) ([]graphql.InputObjectResolver, error) {
	var res []graphql.InputObjectResolver
	for _, msg := range file.File.Messages {
		msgCfg, err := file.Config.MessageConfig(msg.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message '%s' config", msg.Name)
		}

		if msgCfg.UnwrapField {
			resolver, err := g.inputUnwrappedMessagesResolver(file, msg)
			if err != nil {
				return nil, err
			}

			res = append(res, *resolver)

			continue
		}

		var oneOffs []graphql.InputObjectResolverOneOf
		for _, oneOf := range msg.OneOffs {
			var fields []graphql.InputObjectResolverOneOfField
			for _, fld := range oneOf.Fields {
				fldTypeFile, err := g.parsedFile(fld.Type.File())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve message '%s' field '%s' type parsed file", dotedTypeName(msg.TypeName), fld.Name)
				}
				resolver, withErr, _, err := g.TypeValueResolver(fldTypeFile, fld.Type, "", false)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get type value resolver")
				}
				fields = append(fields, graphql.InputObjectResolverOneOfField{
					GraphQLInputFieldName: fld.Name,
					ValueResolver:         resolver,
					ResolverWithError:     withErr,
					AssigningWrapper:      g.oneOfValueAssigningWrapper(file, msg, fld),
				})
			}
			oneOffs = append(oneOffs, graphql.InputObjectResolverOneOf{
				OutputFieldName: camelCase(oneOf.Name),
				Fields:          fields,
			})
		}
		var fields []graphql.InputObjectResolverField
		for _, fld := range msg.NormalFields {
			fldTypeFile, err := g.parsedFile(fld.Type.File())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message '%s' field '%s' type parsed file", dotedTypeName(msg.TypeName), fld.Name)
			}
			fldCfg := msgCfg.Fields[fld.Name]
			resolver, withErr, fromArgs, err := g.TypeValueResolver(fldTypeFile, fld.Type, fldCfg.ContextKey, fld.Optional || fld.Required)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get type value resolver")
			}
			goType, err := g.goTypeByParserType(fld.Type)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get go type by parser type")
			}
			if fld.Repeated {
				gt := goType
				goType = graphql.GoType{
					Kind:     reflect.Slice,
					ElemType: &gt,
				}
			}
			fields = append(fields, graphql.InputObjectResolverField{
				GraphQLInputFieldName: fld.Name,
				OutputFieldName:       camelCase(fld.Name),
				ValueResolver:         resolver,
				ResolverWithError:     withErr,
				GoType:                goType,
				IsFromArgs:            fromArgs,
			})
		}
		for _, fld := range msg.MapFields {
			valueTypeParsedFile, err := g.parsedFile(fld.Map.File())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message '%s' parsed file", dotedTypeName(msg.TypeName))
			}
			fldCfg := msgCfg.Fields[fld.Name]
			valueResolver, withErr, fromArgs, err := g.TypeValueResolver(valueTypeParsedFile, fld.Map, fldCfg.ContextKey, false)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get message '%s' map field '%s' value resolver", msg.Name, fld.Name)
			}
			fields = append(fields, graphql.InputObjectResolverField{
				GraphQLInputFieldName: fld.Name,
				OutputFieldName:       camelCase(fld.Name),
				ValueResolver:         valueResolver,
				ResolverWithError:     withErr,
				IsFromArgs:            fromArgs,
			})
		}
		msgGoType, err := g.goTypeByParserType(msg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve message go type")
		}
		res = append(res, graphql.InputObjectResolver{
			FunctionName: g.inputMessageResolverName(file, msg),
			OutputGoType: msgGoType,
			OneOfFields:  oneOffs,
			Fields:       fields,
		})

	}

	return res, nil
}
