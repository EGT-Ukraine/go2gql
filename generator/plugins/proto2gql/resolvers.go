package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) TypeOutputTypeResolver(typeFile *parsedFile, typ parser.Type) (graphql.TypeResolver, error) {
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
		res, err := g.outputMessageTypeResolver(typeFile, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get message type resolver")
		}
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
	return nil, errors.New("not implemented " + typ.String())
}
func (g *Proto2GraphQL) TypeInputTypeResolver(typeFile *parsedFile, typ parser.Type) (graphql.TypeResolver, error) {
	switch pType := typ.(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[pType.ScalarName]
		if !ok {
			return nil, errors.Errorf("unimplemented scalar type: %s", typ.(*parser.Scalar).ScalarName)
		}
		return resolver, nil
	case *parser.Message:
		res, err := g.inputMessageTypeResolver(typeFile, pType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get message type resolver")
		}
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
func (g *Proto2GraphQL) TypeValueResolver(typeFile *parsedFile, typ parser.Type, ctxKey string) (_ graphql.ValueResolver, withErr, fromArgs bool, err error) {
	if ctxKey != "" {
		goType, err := g.goTypeByParserType(typ)
		if err != nil {
			return nil, false, false, errors.Wrap(err, "failed to resolve go type")
		}
		return func(arg string, ctx graphql.BodyContext) string {
			return `ctx.Value("` + ctxKey + `").(` + goType.String(ctx.Importer) + `)`
		}, false, false, nil
	}
	switch pType := typ.(type) {
	case *parser.Scalar:
		gt, ok := goTypesScalars[pType.ScalarName]
		if !ok {
			panic("unknown scalar: " + pType.ScalarName)
		}
		return func(arg string, ctx graphql.BodyContext) string {
			return arg + ".(" + gt.Kind.String() + ")"
		}, false, true, nil
	case *parser.Enum:
		return func(arg string, ctx graphql.BodyContext) string {
			return ctx.Importer.Prefix(typeFile.GRPCSourcesPkg) + snakeCamelCaseSlice(pType.TypeName) + "(" + arg + ".(int))"
		}, false, true, nil
	case *parser.Message:
		return func(arg string, ctx graphql.BodyContext) string {

			if ctx.TracerEnabled {
				return ctx.Importer.Prefix(typeFile.OutputPkg) + g.inputMessageResolverName(typeFile, pType) + "(tr, " + ctx.Importer.New(graphql.OpentracingPkgPath) + ".ContextWithSpan(ctx,span), " + arg + ")"
			} else {
				return ctx.Importer.Prefix(typeFile.OutputPkg) + g.inputMessageResolverName(typeFile, pType) + "(ctx, " + arg + ")"
			}
		}, true, true, nil

	case *parser.Map:
		return func(arg string, ctx graphql.BodyContext) string {
			if ctx.TracerEnabled {
				return ctx.Importer.Prefix(typeFile.OutputPkg) + g.mapResolverFunctionName(typeFile, pType) + "(tr, " + ctx.Importer.New(graphql.OpentracingPkgPath) + ".ContextWithSpan(ctx,span), " + arg + ")"
			} else {
				return ctx.Importer.Prefix(typeFile.OutputPkg) + g.mapResolverFunctionName(typeFile, pType) + "(ctx, " + arg + ")"
			}
		}, true, true, nil
	}

	return func(arg string, ctx graphql.BodyContext) string {
		return arg + "// not implemented"
	}, false, true, nil

}

func (g *Proto2GraphQL) FieldOutputValueResolver(fieldFile *parsedFile, fieldName string, fieldRepeated bool, fieldType parser.Type) (_ graphql.ValueResolver, err error) {
	switch ft := fieldType.(type) {
	case *parser.Scalar:
		return graphql.IdentAccessValueResolver(camelCase(fieldName)), nil
	case *parser.Message:
		return graphql.IdentAccessValueResolver(camelCase(fieldName)), nil
	case *parser.Map:
		goKeyTyp, err := g.goTypeByParserType(ft.KeyType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve field key go type")
		}
		goValueTyp, err := g.goTypeByParserType(ft.ValueType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve field value go type")
		}
		return func(arg string, ctx graphql.BodyContext) string {
			return "func(arg map[" + goKeyTyp.String(ctx.Importer) + "]" + goValueTyp.String(ctx.Importer) + ") []map[string]interface{} {" +
				"\n  	res := make([]int, len(arg))" +
				"\n 	for i, val := range arg {" +
				"\n 		res[i] = int(val)" +
				"\n		}" +
				"\n 	return res" +
				"\n	}(" + arg + ".Get" + camelCase(fieldName) + "())"
		}, nil
	case *parser.Enum:
		if fieldRepeated {
			goTyp, err := g.goTypeByParserType(fieldType)
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
		} else {
			return func(arg string, ctx graphql.BodyContext) string {
				return "int(" + arg + ".Get" + camelCase(fieldName) + "())"
			}, nil
		}
	}
	return func(arg string, ctx graphql.BodyContext) string {
		return arg + "// not implemented"
	}, nil
}
