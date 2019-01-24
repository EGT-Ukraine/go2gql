package swagger2gql

import (
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) mapResolverFunctionName(file *parsedFile, mp *parser.Map) string {
	return "Resolve" + p.mapInputObjectVariable(file, mp)
}

func (p *Plugin) fileInputMapResolvers(file *parsedFile) ([]graphql.MapInputObjectResolver, error) {
	var res []graphql.MapInputObjectResolver
	handledObjects := map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
		case *parser.Map:
			valueGoType, err := p.goTypeByParserType(file, t.ElemType, false)
			if err != nil {
				return errors.Wrap(err, "failed to resolve map value go type")
			}
			valueResolver, valueWithErr, _, err := p.TypeValueResolver(file, t.ElemType, false, "")

			if err != nil {
				return err
			}

			res = append(res, graphql.MapInputObjectResolver{
				FunctionName: p.mapResolverFunctionName(file, t),
				KeyGoType: graphql.GoType{
					Kind: reflect.String,
				},
				ValueGoType: valueGoType,
				KeyResolver: func(arg string, ctx graphql.BodyContext) string {
					return arg + ".(string)"
				},
				KeyResolverWithError:   false,
				ValueResolver:          valueResolver,
				ValueResolverWithError: valueWithErr,
			})
		case *parser.Object:
			if _, handled := handledObjects[t]; handled {
				return nil
			}
			handledObjects[t] = struct{}{}
			for _, property := range t.Properties {
				if err := handleType(property.Type); err != nil {
					return errors.Wrapf(err, "failed to handle object property %s type", property.Name)
				}
			}
		case *parser.Array:
			return handleType(t.ElemType)
		}
		return nil
	}
	for _, tag := range file.File.Tags {
		for _, method := range tag.Methods {
			for _, param := range method.Parameters {
				if err := handleType(param.Type); err != nil {
					return nil, errors.Wrapf(err, "failed to handle %s method %s parameter", method.OperationID, param.Name)
				}

			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].FunctionName > res[j].FunctionName
	})
	return res, nil
}
