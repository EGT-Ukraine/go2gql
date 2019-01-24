package swagger2gql

import (
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/names"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

func (p *Plugin) outputObjectGQLName(messageFile *parsedFile, obj *parser.Object) string {
	return messageFile.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, "__"))
}
func (p *Plugin) outputObjectVariable(messageFile *parsedFile, obj *parser.Object) string {
	return messageFile.Config.GetGQLMessagePrefix() + pascalize(strings.Join(obj.Route, ""))
}

func (p *Plugin) outputMessageTypeResolver(messageFile *parsedFile, obj *parser.Object) graphql.TypeResolver {
	if len(obj.Properties) == 0 {
		return graphql.GqlNoDataTypeResolver
	}
	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(messageFile.OutputPkg) + p.outputObjectVariable(messageFile, obj)
	}
}

func (p *Plugin) fileOutputMessages(file *parsedFile) ([]graphql.OutputObject, error) {
	var res []graphql.OutputObject
	handledObjects := map[parser.Type]struct{}{}
	var handleType func(typ parser.Type) error
	handleType = func(typ parser.Type) error {
		switch t := typ.(type) {
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
			goTyp, err := p.goTypeByParserType(file, t, false)
			if err != nil {
				return errors.Wrap(err, "failed to resolve object go type")
			}
			var fields []graphql.ObjectField
			var mapFields []graphql.ObjectField
			for _, prop := range t.Properties {
				tr, err := p.TypeOutputTypeResolver(file, prop.Type, false)
				if err != nil {
					return errors.Wrap(err, "failed to resolve property output type resolver")
				}
				valueResolver := graphql.IdentAccessValueResolver(pascalize(prop.Name))
				if typ == parser.ObjDateTime {
					switch prop.Name {
					case "seconds":
						valueResolver = func(arg string, ctx graphql.BodyContext) string {
							return `(time.Time)(` + arg + `).Unix()`
						}
					case "nanos":
						valueResolver = func(arg string, ctx graphql.BodyContext) string {
							return `int32((time.Time)(` + arg + `).Nanosecond())`
						}
					}

				}

				propGoType, err := p.goTypeByParserType(file, prop.Type, false)
				if err != nil {
					return errors.Wrap(err, "failed to resolve property go type")
				}

				propObj := graphql.ObjectField{
					Name:          names.FilterNotSupportedFieldNameCharacters(prop.Name),
					QuotedComment: strconv.Quote(prop.Description),
					Type:          tr,
					Value:         valueResolver,
					NeedCast:      false,
					GoType:        propGoType,
				}
				if prop.Type.Kind() == parser.KindMap {
					mapFields = append(mapFields, propObj)

				} else {
					fields = append(fields, propObj)
				}
			}
			sort.Slice(fields, func(i, j int) bool {
				return fields[i].Name > fields[j].Name
			})

			objectName := p.outputObjectGQLName(file, t)
			objectConfig, err := file.Config.ObjectConfig(objectName)

			if err != nil {
				return errors.Wrap(err, "failed to get object config "+objectName)
			}

			dataLoaderFields, err := p.dataLoaderFields(objectConfig.DataLoaders, t)
			if err != nil {
				return errors.Wrapf(err, "failed to resolve output object %s data loaders", objectName)
			}

			res = append(res, graphql.OutputObject{
				VariableName:     p.outputObjectVariable(file, t),
				GraphQLName:      p.outputObjectGQLName(file, t),
				GoType:           goTyp,
				Fields:           fields,
				MapFields:        mapFields,
				DataLoaderFields: dataLoaderFields,
			})
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
	sort.Slice(res, func(i, j int) bool {
		return res[i].VariableName > res[j].VariableName
	})
	return res, nil
}

func (p *Plugin) dataLoaderFields(configs []dataloader.FieldConfig, object *parser.Object) ([]*graphql.DataLoaderField, error) {
	var fields []*graphql.DataLoaderField

	for _, cfg := range configs {
		prop := object.GetPropertyByName(cfg.KeyFieldName)

		if prop == nil {
			return nil, errors.Errorf("Can't find property %s for dataloader", cfg.KeyFieldName)
		}

		field := &graphql.DataLoaderField{
			Name:                         cfg.FieldName,
			ParentKeyFieldName:           cfg.KeyFieldName,
			KeyFieldSlice:                prop.Type.Kind() == parser.KindArray,
			NormalizedParentKeyFieldName: pascalize(cfg.KeyFieldName),
			DataLoaderName:               cfg.DataLoader,
		}

		fields = append(fields, field)
	}

	return fields, nil
}
