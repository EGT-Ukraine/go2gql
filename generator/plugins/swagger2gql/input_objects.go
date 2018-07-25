package swagger2gql

import (
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/names"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) inputObjectGQLName(file *parsedFile, obj *parser.Object) string {
	return file.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, "__")) + "Input"
}
func (p *Plugin) inputObjectVariable(msgFile *parsedFile, obj *parser.Object) string {
	return msgFile.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, "")) + "Input"
}
func (p *Plugin) methodParamsInputObjectVariable(file *parsedFile, method parser.Method) string {
	return file.Config.GetGQLMessagePrefix() + pascalize(method.OperationID+"Params") + "Input"
}
func (p *Plugin) methodParamsInputObjectGQLName(file *parsedFile, method parser.Method) string {
	return file.Config.GetGQLMessagePrefix() + pascalize(method.OperationID+"Params") + "Input"
}

func (p *Plugin) inputObjectTypeResolver(msgFile *parsedFile, obj *parser.Object) graphql.TypeResolver {
	if len(obj.Properties) == 0 {
		return graphql.GqlNoDataTypeResolver
	}

	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(msgFile.OutputPkg) + p.inputObjectVariable(msgFile, obj)
	}
}

func (p *Plugin) fileInputObjects(file *parsedFile) ([]graphql.InputObject, error) {
	var res []graphql.InputObject
	var handledObjects = map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
		case *parser.Object:
			if _, handled := handledObjects[typ]; handled {
				return nil
			}
			handledObjects[typ] = struct{}{}
			var fields []graphql.ObjectField
			for _, property := range t.Properties {
				if err := handleType(property.Type); err != nil {
					return err
				}
				typeResolver, err := p.TypeInputTypeResolver(file, property.Type)
				if err != nil {
					return errors.Wrap(err, "failed to get input type resolver")
				}
				if property.Required {
					typeResolver = graphql.GqlNonNullTypeResolver(typeResolver)
				}
				fields = append(fields, graphql.ObjectField{
					Name:     names.FilterNotSupportedFieldNameCharacters(property.Name),
					Type:     typeResolver,
					NeedCast: false,
				})
			}
			sort.Slice(fields, func(i, j int) bool {
				return fields[i].Name > fields[j].Name
			})
			res = append(res, graphql.InputObject{
				VariableName: p.inputObjectVariable(file, t),
				GraphQLName:  p.inputObjectGQLName(file, t),
				Fields:       fields,
			})

		case *parser.Array:
			return handleType(t.ElemType)
		}
		return nil
	}
	for _, tag := range file.File.Tags {
		for _, method := range tag.Methods {
			for _, parameter := range method.Parameters {
				err := handleType(parameter.Type)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to handle method %s parameter %s", method.OperationID, parameter.Name)
				}
			}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].VariableName > res[j].VariableName
	})
	return res, nil
}
