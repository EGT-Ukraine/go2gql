package swagger2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

var scalarsResolvers = map[parser.Kind]graphql.TypeResolver{
	parser.KindBoolean:  graphql.GqlBoolTypeResolver,
	parser.KindFloat64:  graphql.GqlFloat64TypeResolver,
	parser.KindFloat32:  graphql.GqlFloat32TypeResolver,
	parser.KindInt64:    graphql.GqlInt64TypeResolver,
	parser.KindInt32:    graphql.GqlInt32TypeResolver,
	parser.KindString:   graphql.GqlStringTypeResolver,
	parser.KindNull:     graphql.GqlNoDataTypeResolver,
	parser.KindFile:     graphql.GqlMultipartFileTypeResolver,
	parser.KindDateTime: graphql.GqlStringTypeResolver,
}

func (p *Plugin) TypeOutputTypeResolver(typeFile *parsedFile, typ parser.Type, required bool) (graphql.TypeResolver, error) {
	var res graphql.TypeResolver
	switch t := typ.(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[typ.Kind()]
		if !ok {
			return nil, errors.Errorf(": %s", typ.Kind())
		}
		res = resolver
	case *parser.Object:
		msgResolver, err := p.outputMessageTypeResolver(typeFile, t)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get message type resolver")
		}
		res = msgResolver
	case *parser.Array:
		elemResolver, err := p.TypeOutputTypeResolver(typeFile, t.ElemType, true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get array element type resolver")
		}
		res = graphql.GqlListTypeResolver(elemResolver)
	case *parser.Map:
		res = func(ctx graphql.BodyContext) string {
			return p.mapOutputObjectVariable(typeFile, t)
		}
		res = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(res))
	default:
		return nil, errors.Errorf("not implemented %v", typ.Kind())
	}
	if required {
		res = graphql.GqlNonNullTypeResolver(res)
	}
	return res, nil
}
func (p *Plugin) TypeInputTypeResolver(typeFile *parsedFile, typ parser.Type) (graphql.TypeResolver, error) {
	switch t := typ.(type) {
	case *parser.Scalar:
		resolver, ok := scalarsResolvers[t.Kind()]
		if !ok {
			return nil, errors.Errorf("unimplemented scalar type: %s", t.Kind())
		}
		return resolver, nil
	case *parser.Object:
		return p.inputObjectTypeResolver(typeFile, t), nil
	case *parser.Array:
		elemResolver, err := p.TypeInputTypeResolver(typeFile, t.ElemType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get array element type resolver")
		}
		return graphql.GqlListTypeResolver(elemResolver), nil
	case *parser.Map:
		res := func(ctx graphql.BodyContext) string {
			return p.mapInputObjectVariable(typeFile, t)
		}
		return graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(res)), nil
	}
	return nil, errors.New("not implemented " + typ.Kind().String())
}
func (p *Plugin) TypeValueResolver(file *parsedFile, typ parser.Type, required bool, ctxKey string) (_ graphql.ValueResolver, withErr, fromArgs bool, err error) {
	if ctxKey != "" {
		goType, err := p.goTypeByParserType(file, typ, true)
		if err != nil {
			return nil, false, false, errors.Wrap(err, "failed to resolve go type")
		}
		return func(arg string, ctx graphql.BodyContext) string {
			valueType := goType.String(ctx.Importer)

			if !required {
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
			}

			return `func() (*` + valueType + `, error) {
							contextValue := ctx.Value("` + ctxKey + `")

							if contextValue == nil {
								return nil, errors.New("Can't find key '` + ctxKey + `' in context")
							}

							val, ok := contextValue.(` + valueType + `)

							if !ok {
								return nil, errors.New("Incompatible '` + ctxKey + `' key type in context. Expected ` + valueType + `")
							}

							return &val, nil
						}()`
		}, true, false, nil
	}
	switch t := typ.(type) {
	case *parser.Scalar:
		if t.Kind() == parser.KindFile {
			return func(arg string, ctx graphql.BodyContext) string {
				return "(" + arg + ").(*" + ctx.Importer.Prefix(graphql.MultipartFilePkgPath) + "MultipartFile)"
			}, false, true, nil
		}
		goTyp, ok := scalarsGoTypesNames[typ.Kind()]
		if !ok {
			return nil, false, false, errors.Errorf("scalar %s is not implemented", typ.Kind())
		}
		return func(arg string, ctx graphql.BodyContext) string {
			if !required {
				return arg + ".(" + goTyp + ")"
			}
			return "func(arg interface{}) *" + goTyp + "{\n" +
				"val := arg.(" + goTyp + ")\n" +
				"return &val\n" +
				"}(" + arg + ")"
		}, false, true, nil
	case *parser.Object:
		if t == parser.ObjDateTime {
			return func(arg string, ctx graphql.BodyContext) string {
				if required {
					res, err := p.renderPtrDatetimeResolver(arg, ctx)
					if err != nil {
						panic(errors.Wrap(err, "failed to render ptr datetime resolver"))
					}
					return res
				} else {
					res, err := p.renderDatetimeValueResolverTemplate(arg, ctx)
					if err != nil {
						panic(errors.Wrap(err, "failed to render ptr datetime resolver"))
					}
					return res
				}
			}, true, true, nil
		}
		return graphql.ResolverCall(file.OutputPkg, "Resolve"+snakeCamelCaseSlice(t.Route)), true, true, nil

	case *parser.Array:
		elemResolver, elemResolverWithErr, _, err := p.TypeValueResolver(file, t.ElemType, false, "")
		if err != nil {
			return nil, false, false, errors.Wrap(err, "failed to get array element type value resolver")
		}
		goTyp, err := p.goTypeByParserType(file, typ, true)
		if err != nil {
			return nil, false, false, errors.Wrap(err, "failed to resolve go type by parser type")
		}
		return func(arg string, ctx graphql.BodyContext) string {
			res, err := p.renderArrayValueResolver(arg, goTyp, ctx, elemResolver, elemResolverWithErr)
			if err != nil {
				panic(err)
			}
			return res
		}, true, true, nil
	case *parser.Map:
		return graphql.ResolverCall(file.OutputPkg, p.mapResolverFunctionName(file, t)), true, true, nil
	}
	return nil, false, true, errors.Errorf("unknown type: %v", typ.Kind().String())

}
