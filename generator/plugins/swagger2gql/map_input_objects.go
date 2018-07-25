package swagger2gql

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) mapInputObjectGQLName(messageFile *parsedFile, obj *parser.Map) string {
	return messageFile.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, "__")) + "Input"
}
func (p *Plugin) mapInputObjectVariable(messageFile *parsedFile, obj *parser.Map) string {
	return messageFile.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, "")) + "Input"
}

func (p *Plugin) mapInputMessageTypeResolver(messageFile *parsedFile, obj *parser.Map) (graphql.TypeResolver, error) {
	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(messageFile.OutputPkg) + p.mapInputObjectVariable(messageFile, obj)
	}, nil
}
func (p *Plugin) fileMapInputMessages(file *parsedFile) ([]graphql.MapInputObject, error) {
	var res []graphql.MapInputObject
	handledObjects := map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
		case *parser.Map:
			valueType, err := p.TypeInputTypeResolver(file, t.ElemType)
			if err != nil {
				return errors.Wrap(err, "failed to resolve map key type")
			}

			res = append(res, graphql.MapInputObject{
				VariableName:    p.mapInputObjectVariable(file, t),
				GraphQLName:     p.mapInputObjectGQLName(file, t),
				KeyObjectType:   graphql.GqlNonNullTypeResolver(graphql.GqlStringTypeResolver),
				ValueObjectType: valueType,
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
	return res, nil
}
