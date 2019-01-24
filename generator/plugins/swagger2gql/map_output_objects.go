package swagger2gql

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) mapOutputObjectGQLName(messageFile *parsedFile, obj *parser.Map) string {
	return messageFile.Config.GetGQLMessagePrefix() + strings.Join(obj.Route, "__")
}
func (p *Plugin) mapOutputObjectVariable(messageFile *parsedFile, obj *parser.Map) string {
	return messageFile.Config.GetGQLMessagePrefix() + strings.Join(obj.Route, "")
}

func (p *Plugin) fileMapOutputMessages(file *parsedFile) ([]graphql.MapOutputObject, error) {
	var res []graphql.MapOutputObject
	handledObjects := map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
		case *parser.Map:
			valueType, err := p.TypeOutputTypeResolver(file, t.ElemType, true)
			if err != nil {
				return errors.Wrap(err, "failed to resolve map key type")
			}

			res = append(res, graphql.MapOutputObject{
				VariableName:    p.mapOutputObjectVariable(file, t),
				GraphQLName:     p.mapOutputObjectGQLName(file, t),
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
			for _, resp := range method.Responses {
				if err := handleType(resp.ResultType); err != nil {
					return nil, errors.Wrapf(err, "failed to handle %s method %d response", method.OperationID, resp.StatusCode)
				}

			}
		}
	}
	return res, nil
}
