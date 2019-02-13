package proto2gql

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) outputMessageGraphQLName(messageFile *parsedFile, message *parser.Message) string {
	return messageFile.Config.GetGQLMessagePrefix() + camelCaseSlice(message.TypeName)
}

func (g *Proto2GraphQL) outputMessageVariable(messageFile *parsedFile, message *parser.Message) string {
	return messageFile.Config.GetGQLMessagePrefix() + snakeCamelCaseSlice(message.TypeName)
}

func (g *Proto2GraphQL) outputMessageTypeResolver(messageFile *parsedFile, message *parser.Message) graphql.TypeResolver {
	if !message.HaveFields() {
		return graphql.GqlNoDataTypeResolver
	}

	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(messageFile.OutputPkg) + g.outputMessageVariable(messageFile, message)
	}
}

func (g *Proto2GraphQL) outputMessageFields(msgCfg MessageConfig, msg *parser.Message) ([]graphql.ObjectField, error) {
	if msgCfg.UnwrapField {
		fields := msg.GetFields()
		if len(fields) != 1 {
			return nil, errors.New("can't unwrap %s output message because it contains more that 1 field")
		}
		return nil, nil
	}
	var res []graphql.ObjectField
	for _, field := range msg.NormalFields {
		if msgCfg.ErrorField == field.Name {
			continue
		}
		fieldGoType, err := g.goTypeByParserType(field.Type)

		if err != nil {
			return nil, errors.New("failed to resolve property go type")
		}
		fieldMessage, ok := field.Type.(*parser.Message)
		if ok {
			fieldMessageFile, err := g.parsedFile(fieldMessage.File())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message %s parsed file", fieldMessage.Name)
			}
			fieldMessageConfig, err := fieldMessageFile.Config.MessageConfig(fieldMessage.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message %s config", fieldMessage.Name)
			}
			if fieldMessageConfig.UnwrapField {
				object, err := g.outputMessageUnwrappedField(msg, fieldMessage, field)
				if err != nil {
					return nil, errors.Wrap(err, "failed to resolve output message unwrapped field")
				}
				res = append(res, *object)
				continue
			}
		}
		fieldTypeFile, err := g.parsedFile(field.Type.File())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve file type file")
		}

		typeResolver, err := g.TypeOutputGraphQLTypeResolver(fieldTypeFile, field.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare message %s field %s output type resolver", msg.Name, field.Name)
		}

		if field.Repeated {
			typeResolver = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(typeResolver))
		}
		valueResolver, err := g.FieldOutputValueResolver(msg, field.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare message %s field %s output value resolver", msg.Name, field.Name)
		}

		res = append(res, graphql.ObjectField{
			Name:          field.Name,
			QuotedComment: field.QuotedComment,
			Type:          typeResolver,
			GoType:        fieldGoType,
			Value:         valueResolver,
		})
	}
	for _, of := range msg.OneOffs {
		for _, field := range of.Fields {
			if msgCfg.ErrorField == field.Name {
				continue
			}
			fieldTypeFile, err := g.parsedFile(field.Type.File())
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve file type file")
			}
			typeResolver, err := g.TypeOutputGraphQLTypeResolver(fieldTypeFile, field.Type)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to prepare message %s field %s output type resolver", msg.Name, field.Name)
			}
			res = append(res, graphql.ObjectField{
				Name:          field.Name,
				QuotedComment: field.QuotedComment,
				Type:          typeResolver,
				Value:         graphql.IdentAccessValueResolver("Get" + camelCase(field.Name) + "()"),
			})
		}
	}

	return res, nil
}

func (g *Proto2GraphQL) outputMessageMapFields(msgCfg MessageConfig, file *parsedFile, msg *parser.Message) ([]graphql.ObjectField, error) {
	var res []graphql.ObjectField
	for _, field := range msg.MapFields {
		if msgCfg.ErrorField == field.Name {
			continue
		}
		typeResolver, err := g.TypeOutputGraphQLTypeResolver(file, field.Map)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare message %s field %s output type resolver", msg.Name, field.Name)
		}
		res = append(res, graphql.ObjectField{
			Name:  field.Name,
			Type:  typeResolver,
			Value: graphql.IdentAccessValueResolver(camelCase(field.Name)),
		})
	}

	return res, nil
}

func (g *Proto2GraphQL) dataLoaderFields(configs []dataloader.FieldConfig, msg *parser.Message) ([]*graphql.DataLoaderField, error) {
	var fields []*graphql.DataLoaderField

	for _, cfg := range configs {
		msgKeyField, ok := msg.GetFieldByName(cfg.KeyFieldName)
		if !ok {
			return nil, errors.Errorf("can't find key field %s for dataloader", cfg.KeyFieldName)
		}

		normalKeyField, ok := msgKeyField.(*parser.NormalField)
		if !ok {
			return nil, errors.Errorf("only normal fields(not maps) can be keys for dataloaders")
		}

		field := &graphql.DataLoaderField{
			Name:                         cfg.FieldName,
			NormalizedParentKeyFieldName: camelCase(cfg.KeyFieldName),
			ParentKeyFieldName:           cfg.KeyFieldName,
			KeyFieldSlice:                normalKeyField.Repeated,
			DataLoaderName:               cfg.DataLoader,
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func (g *Proto2GraphQL) fileOutputMessages(file *parsedFile) ([]graphql.OutputObject, error) {
	var res []graphql.OutputObject
	for _, msg := range file.File.Messages {
		cfg, err := file.Config.MessageConfig(msg.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s config", msg.Name)
		}
		fields, err := g.outputMessageFields(cfg, msg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s fields", msg.Name)
		}
		mapFields, err := g.outputMessageMapFields(cfg, file, msg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s map fields", msg.Name)
		}

		dataLoaderFields, err := g.dataLoaderFields(cfg.DataLoaders, msg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s data loaders", msg.Name)
		}

		if len(fields)+len(mapFields)+len(dataLoaderFields) == 0 {
			// It's a NoData scalar
			continue
		}

		res = append(res, graphql.OutputObject{
			VariableName: g.outputMessageVariable(file, msg),
			GraphQLName:  g.outputMessageGraphQLName(file, msg),
			Fields:       fields,
			MapFields:    mapFields,
			GoType: graphql.GoType{
				Kind: reflect.Struct,
				Name: snakeCamelCaseSlice(msg.TypeName),
				Pkg:  file.GRPCSourcesPkg,
			},
			DataLoaderFields: dataLoaderFields,
		})
	}

	return res, nil
}

func (g *Proto2GraphQL) outputMessageUnwrappedField(msg, fieldMessage *parser.Message, field parser.Field) (*graphql.ObjectField, error) {
	unwrappedField := fieldMessage.NormalFields[0]
	unwrappedFieldTypeFile, err := g.parsedFile(unwrappedField.Type.File())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve unwrappedField %s message file", unwrappedField.Name)
	}
	outProtoType := unwrappedField.Type

	typeResolver, err := g.TypeOutputGraphQLTypeResolver(unwrappedFieldTypeFile, outProtoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get type output type resolver")
	}
	unwrappedFieldGoType, err := g.goTypeByParserType(outProtoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve go type of unwrapped unwrappedField")
	}
	valueResolver, err := g.FieldOutputValueResolver(msg, field.GetName())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare message %s unwrappedField %s output value resolver", msg.Name, unwrappedField.Name)
	}
	if field.IsRepeated() {
		typeResolver = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(typeResolver))
	}
	return &graphql.ObjectField{
		Name:          field.GetName(),
		QuotedComment: unwrappedField.QuotedComment,
		Type:          typeResolver,
		GoType:        unwrappedFieldGoType,
		Value:         valueResolver,
	}, nil
}
