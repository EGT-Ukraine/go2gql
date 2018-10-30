package swagger2gql

import (
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/names"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) objectResolverFunctionName(file *parsedFile, obj *parser.Object) string {
	return "Resolve" + snakeCamelCaseSlice(obj.Route)
}
func (p *Plugin) methodParametersObjectResolverFuncName(file *parsedFile, method parser.Method) string {
	return "Resolve" + pascalize(method.OperationID) + "Params"
}

func (p *Plugin) methodParametersInputObjectResolver(file *parsedFile, tag string, method parser.Method) (*graphql.InputObjectResolver, error) {
	var fields []graphql.InputObjectResolverField
	gqlName := p.methodParamsInputObjectGQLName(file, method)

	for _, param := range method.Parameters {
		// goTyp, err := p.goTypeByParserType(file, param.Type, true)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "failed to resolve parameter go type")
		// }
		paramGqlName := names.FilterNotSupportedFieldNameCharacters(param.Name)
		paramCfg, err := file.Config.FieldConfig(gqlName, param.Name)

		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve property config")
		}

		valueResolver, withErr, fromArgs, err := p.TypeValueResolver(file, param.Type, !param.Required, paramCfg.ContextKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parameter value resolver")
		}
		fields = append(fields, graphql.InputObjectResolverField{
			OutputFieldName:       pascalize(param.Name),
			GraphQLInputFieldName: paramGqlName,
			ValueResolver:         valueResolver,
			ResolverWithError:     withErr,
			IsFromArgs:            fromArgs,
		})
	}

	return &graphql.InputObjectResolver{
		FunctionName: p.methodParametersObjectResolverFuncName(file, method),
		Fields:       fields,
		OutputGoType: graphql.GoType{
			Kind: reflect.Ptr,
			ElemType: &graphql.GoType{
				Kind: reflect.Struct,
				Name: pascalize(method.OperationID) + "Params",
				Pkg:  file.Config.Tags[tag].ClientGoPackage,
			},
		},
	}, nil
}

func (p *Plugin) fileInputMessagesResolvers(file *parsedFile) ([]graphql.InputObjectResolver, error) {
	var res []graphql.InputObjectResolver
	var handledObjects = map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
		case *parser.Array:
			return handleType(t.ElemType)
		case *parser.Object:
			if t == parser.ObjDateTime {
				return nil
			}
			var fields []graphql.InputObjectResolverField
			if _, handled := handledObjects[t]; handled {
				return nil
			}
			gqlObjName := p.inputObjectGQLName(file, t)
			handledObjects[t] = struct{}{}
			for _, property := range t.Properties {
				gqlName := names.FilterNotSupportedFieldNameCharacters(property.Name)

				paramCfg, err := file.Config.FieldConfig(gqlObjName, property.Name)
				if err != nil {
					return errors.Wrapf(err, "failed to resolve property %s config", property.Name)
				}

				err = handleType(property.Type)
				if err != nil {
					return errors.Wrapf(err, "failed to resolve property %s objects resolvers", property.Name)
				}
				valueResolver, withErr, fromArgs, err := p.TypeValueResolver(file, property.Type, property.Required, paramCfg.ContextKey)
				if err != nil {
					return errors.Wrap(err, "failed to get property value resolver")
				}
				fields = append(fields, graphql.InputObjectResolverField{
					GraphQLInputFieldName: gqlName,
					OutputFieldName:       pascalize(property.Name),
					ValueResolver:         valueResolver,
					ResolverWithError:     withErr,
					GoType: graphql.GoType{
						Kind:   reflect.Uint,
						Scalar: true,
					},
					IsFromArgs: fromArgs,
				})
			}
			resGoType, err := p.goTypeByParserType(file, t, true)
			if err != nil {
				return errors.Wrap(err, "failed to resolve object go type")
			}
			sort.Slice(fields, func(i, j int) bool {
				return fields[i].GraphQLInputFieldName > fields[j].GraphQLInputFieldName
			})
			res = append(res, graphql.InputObjectResolver{
				FunctionName: p.objectResolverFunctionName(file, t),
				Fields:       fields,
				OutputGoType: resGoType,
			})

		}
		return nil
	}
	for _, tag := range file.File.Tags {
		_, ok := file.Config.Tags[tag.Name]

		if !ok {
			continue
		}

		for _, method := range tag.Methods {
			paramsResolver, err := p.methodParametersInputObjectResolver(file, tag.Name, method)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get method partameters input object resolver")
			}
			res = append(res, *paramsResolver)
			for _, parameter := range method.Parameters {
				err := handleType(parameter.Type)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to handle type %v", parameter.Type.Kind())
				}
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].FunctionName > res[j].FunctionName
	})
	return res, nil
}
