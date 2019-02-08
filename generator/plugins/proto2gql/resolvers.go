package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) FieldOutputGraphQLTypeResolver(message *parser.Message, fieldName string) (graphql.TypeResolver, error) {
	field, ok := message.GetFieldByName(fieldName)
	if !ok {
		return nil, errors.Errorf("can't find field %s of message %s", fieldName, message.Name)
	}
	fieldTypeFile, err := g.parsedFile(field.GetType().File())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to field %s type parsed filed", field.GetName())
	}

	var result graphql.TypeResolver

	switch pType := field.GetType().(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[pType.ScalarName]
		if !ok {
			return nil, errors.Errorf("unimplemented scalar type: %s", pType.ScalarName)
		}
		result = resolver

	case *parser.Message:
		messageConfig, err := fieldTypeFile.Config.MessageConfig(pType.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s config", pType.Name)
		}

		if !pType.HaveFieldsExcept(messageConfig.ErrorField) {
			return graphql.GqlNoDataTypeResolver, nil
		}
		if messageConfig.UnwrapField {
			messageFields := pType.GetFields()
			if len(messageFields) != 1 {
				return nil, errors.Errorf("unwrapped message %s should have one field", pType.Name)
			}
			unwrappedField := messageFields[0]
			result, err = g.FieldOutputGraphQLTypeResolver(pType, unwrappedField.GetName())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to make resolver for %s message path %s", pType.Name, unwrappedField.GetName())
			}
		} else {
			result = g.outputMessageTypeResolver(fieldTypeFile, pType)
		}

	case *parser.Enum:
		file, err := g.parsedFile(pType.File())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve type parsed file")
		}
		res, err := g.enumTypeResolver(file, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get enum type resolver")
		}
		return res, nil

	case *parser.Map:
		return g.outputObjectMapFieldTypeResolver(fieldTypeFile, pType)

	default:
		return nil, errors.Errorf("not implemented %v", field.GetType())
	}
	if field.IsRepeated() {
		result = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(result))
	}

	return result, nil
}

func (g *Proto2GraphQL) TypeOutputGraphQLTypeResolver(typeFile *parsedFile, typ parser.Type) (graphql.TypeResolver, error) {
	switch pType := typ.(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[pType.ScalarName]
		if !ok {
			return nil, errors.Errorf("unimplemented scalar type: %s", pType.ScalarName)
		}
		return resolver, nil
	case *parser.Message:
		msgCfg, err := typeFile.Config.MessageConfig(pType.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s config", pType.Name)
		}
		if !pType.HaveFieldsExcept(msgCfg.ErrorField) {
			return graphql.GqlNoDataTypeResolver, nil
		}
		res := g.outputMessageTypeResolver(typeFile, pType)

		return res, nil
	case *parser.Enum:
		file, err := g.parsedFile(pType.File())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve type parsed file")
		}
		res, err := g.enumTypeResolver(file, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get enum type resolver")
		}
		return res, nil
	case *parser.Map:
		return g.outputObjectMapFieldTypeResolver(typeFile, pType)
	}
	return nil, errors.Errorf("not implemented %v", typ)
}

func (g *Proto2GraphQL) TypeInputGraphQLTypeResolver(typeFile *parsedFile, typ parser.Type) (graphql.TypeResolver, error) {
	switch pType := typ.(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[pType.ScalarName]
		if !ok {
			return nil, errors.Errorf("unimplemented scalar type: %s", typ.(*parser.Scalar).ScalarName)
		}
		return resolver, nil
	case *parser.Message:
		res := g.inputMessageTypeResolver(typeFile, pType)

		return res, nil
	case *parser.Enum:
		res, err := g.enumTypeResolver(typeFile, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get enum type resolver")
		}
		return res, nil
	case *parser.Map:
		return g.inputObjectMapFieldTypeResolver(typeFile, pType)

	}
	return nil, errors.New("not implemented " + typ.String())
}

func (g *Proto2GraphQL) TypeValueResolver(typeFile *parsedFile, typ parser.Type, ctxKey string, ptr bool) (_ graphql.ValueResolver, withErr, fromArgs bool, err error) {
	if ctxKey != "" {
		goType, err := g.goTypeByParserType(typ)
		if err != nil {
			return nil, false, false, errors.Wrap(err, "failed to resolve go type")
		}
		return func(arg string, ctx graphql.BodyContext) string {
			valueType := goType.String(ctx.Importer)

			return `func() (val ` + valueType + `, err error) {
							contextValue := ctx.Value("` + ctxKey + `")

							if contextValue == nil {
								err = errors.New("Can't find key '` + ctxKey + `' in context")
								return
							}

							val, ok := contextValue.(` + valueType + `)

							if !ok {
								err = errors.New("Incompatible '` + ctxKey + `' key type in context. Expected ` + valueType + `")
								return
							}

							return
						}()`
		}, true, false, nil
	}
	switch pType := typ.(type) {
	case *parser.Scalar:
		gt, ok := goTypesScalars[pType.ScalarName]
		if !ok {
			panic("unknown scalar: " + pType.ScalarName)
		}
		return func(arg string, ctx graphql.BodyContext) string {
			goTyp := gt.String(ctx.Importer)
			resolver := arg + ".(" + gt.String(ctx.Importer) + ")"
			if ptr && gt.Kind != graphql.KindBytes {
				resolver = optionalValueResolver(goTyp, resolver, arg)
			}
			return resolver
		}, false, true, nil
	case *parser.Enum:
		return func(arg string, ctx graphql.BodyContext) string {
			goTyp := ctx.Importer.Prefix(typeFile.GRPCSourcesPkg) + snakeCamelCaseSlice(pType.TypeName)
			resolver := goTyp + "(" + arg + ".(int))"
			if ptr {
				resolver = optionalValueResolver(goTyp, resolver, arg)
			}
			return resolver
		}, false, true, nil
	case *parser.Message:
		return func(arg string, ctx graphql.BodyContext) string {
			return ctx.Importer.Prefix(typeFile.OutputPkg) + g.inputMessageResolverName(typeFile, pType) + "(ctx, " + arg + ")"
		}, true, true, nil

	case *parser.Map:
		return func(arg string, ctx graphql.BodyContext) string {
			return ctx.Importer.Prefix(typeFile.OutputPkg) + g.mapResolverFunctionName(typeFile, pType) + "(ctx, " + arg + ")"
		}, true, true, nil
	}

	return func(arg string, ctx graphql.BodyContext) string {
		return arg + "// not implemented"
	}, false, true, nil

}

func (g *Proto2GraphQL) FieldOutputValueResolver(message *parser.Message, fieldName string) (_ graphql.ValueResolver, err error) {
	field, ok := message.GetFieldByName(fieldName)
	if !ok {
		return nil, errors.Errorf("can't find field %s of message %s", fieldName, message.Name)
	}

	var result graphql.ValueResolver
	switch ft := field.GetType().(type) {
	case *parser.Message:
		result = graphql.IdentAccessValueResolver("Get" + camelCase(field.GetName()) + "()")
		messageFile, err := g.parsedFile(ft.File())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s parsed filed", ft.Name)
		}
		messageConfig, err := messageFile.Config.MessageConfig(ft.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s config", ft.Name)
		}
		if messageConfig.UnwrapField {
			childMessageFields := ft.GetFields()
			if len(childMessageFields) != 1 {
				return nil, errors.Errorf("unwrapped message %s should have one field", ft.Name)
			}
			childResolver, err := g.FieldOutputValueResolver(ft, childMessageFields[0].GetName())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to make resolver for %s message path %s", ft.Name, childMessageFields[0].GetName())
			}
			if field.IsRepeated() {
				fieldGoType, err := g.goTypeByParserType(field.GetType())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve field %s go type", field.GetName())
				}

				return repeatedValueResolver(fieldGoType, result, childResolver), nil
			} else {
				return func(arg string, ctx graphql.BodyContext) string {
					return childResolver(result(arg, ctx), ctx)
				}, nil
			}
		}
	case *parser.Scalar:
		result = graphql.IdentAccessValueResolver("Get" + camelCase(field.GetName()) + "()")
	case *parser.Map:
		goKeyTyp, err := g.goTypeByParserType(ft.KeyType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve field key go type")
		}
		goValueTyp, err := g.goTypeByParserType(ft.ValueType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve field value go type")
		}
		result = func(arg string, ctx graphql.BodyContext) string {
			return "func(arg map[" + goKeyTyp.String(ctx.Importer) + "]" + goValueTyp.String(ctx.Importer) + ") []map[string]interface{} {" +
				"\n  	res := make([]int, len(arg))" +
				"\n 	for i, val := range arg {" +
				"\n 		res[i] = int(val)" +
				"\n		}" +
				"\n 	return res" +
				"\n	}(" + arg + ".Get" + camelCase(fieldName) + "())"
		}
	case *parser.Enum:
		if field.IsRepeated() {
			goTyp, err := g.goTypeByParserType(field.GetType())
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve field go type")
			}
			return func(arg string, ctx graphql.BodyContext) string {

				return "func(arg []" + goTyp.String(ctx.Importer) + ") []int {" +
					"\n  	res := make([]int, len(arg))" +
					"\n 	for i, val := range arg {" +
					"\n 		res[i] = int(val)" +
					"\n		}" +
					"\n 	return res" +
					"\n	}(" + arg + ".Get" + camelCase(fieldName) + "())"
			}, nil
		}

		return func(arg string, ctx graphql.BodyContext) string {
			return "int(" + arg + ".Get" + camelCase(fieldName) + "())"
		}, nil
	default:
		return nil, errors.Errorf("can't build output value resolver for field of type %v", field.GetType().Kind())
	}

	return result, nil
}

func optionalValueResolver(goTyp, valueResolver, arg string) string {
	return "func(arg interface{}) *" + goTyp + "{\n" +
		"val := " + valueResolver + "\n" +
		"return &val\n" +
		"}(" + arg + ")"
}

func repeatedValueResolver(fieldGoType graphql.GoType, valueResolver, childValueResolver graphql.ValueResolver) graphql.ValueResolver {
	return func(arg string, ctx graphql.BodyContext) string {
		return `func(values []` + fieldGoType.String(ctx.Importer) + `) (interface{}) {
					var result []interface{}
					for _, value := range values {
						result = append(result, ` + childValueResolver("value", ctx) + `)
					}
					return result
				}(` + valueResolver(arg, ctx) + `)`
	}
}
