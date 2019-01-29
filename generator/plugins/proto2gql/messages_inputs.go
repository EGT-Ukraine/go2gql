package proto2gql

import (
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g *Proto2GraphQL) inputMessageGraphQLName(file *parsedFile, message *parser.Message) string {
	return file.Config.GetGQLMessagePrefix() + camelCaseSlice(message.TypeName) + "Input"
}

func (g *Proto2GraphQL) inputMessageVariable(msgFile *parsedFile, message *parser.Message) string {
	return msgFile.Config.GetGQLMessagePrefix() + snakeCamelCaseSlice(message.TypeName) + "Input"
}

func (g *Proto2GraphQL) inputMessageTypeResolver(msgFile *parsedFile, message *parser.Message) graphql.TypeResolver {
	if !message.HaveFields() {
		return graphql.GqlNoDataTypeResolver
	}

	return func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(msgFile.OutputPkg) + g.inputMessageVariable(msgFile, message)
	}
}

func (g *Proto2GraphQL) inputMessageFieldTypeResolver(file *parsedFile, field *parser.Field) (graphql.TypeResolver, error) {
	resolver, err := g.TypeInputTypeResolver(file, field.Type)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get input type resolver")
	}
	if field.Repeated {
		resolver = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(resolver))
	}

	return resolver, nil
}

func (g *Proto2GraphQL) outputObjectMapFieldTypeResolver(mapFile *parsedFile, mp *parser.Map) (graphql.TypeResolver, error) {
	res := func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(mapFile.OutputPkg) + g.outputMapVariable(mapFile, mp)
	}

	return graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(res)), nil
}

func (g *Proto2GraphQL) inputObjectMapFieldTypeResolver(mapFile *parsedFile, mp *parser.Map) (graphql.TypeResolver, error) {
	res := func(ctx graphql.BodyContext) string {
		return ctx.Importer.Prefix(mapFile.OutputPkg) + g.inputMapVariable(mapFile, mp)
	}

	return graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(res)), nil
}

func (g *Proto2GraphQL) fileInputObjects(file *parsedFile) ([]graphql.InputObject, error) {
	var res []graphql.InputObject
	for _, msg := range file.File.Messages {
		fields, err := g.getMessageFields(file, msg)
		if err != nil {
			return nil, err
		}

		if fields == nil {
			continue
		}

		// TODO: oneof fields
		res = append(res, graphql.InputObject{
			VariableName: g.inputMessageVariable(file, msg),
			GraphQLName:  g.inputMessageGraphQLName(file, msg),
			Fields:       fields,
		})
	}

	return res, nil
}

func (g *Proto2GraphQL) getMessageFields(file *parsedFile, msg *parser.Message) ([]graphql.ObjectField, error) {
	if !msg.HaveFields() {
		return nil, nil
	}

	msgCfg, err := file.Config.MessageConfig(msg.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve message %s config", msg.Name)
	}

	var fields []graphql.ObjectField
	for _, field := range msg.Fields {
		fldCfg := msgCfg.Fields[field.Name]
		if fldCfg.ContextKey != "" {
			continue
		}

		fieldTypeFile, err := g.parsedFile(field.Type.File())
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve value type file")
		}

		fieldTypeMessage, ok := field.Type.(*parser.Message)

		if ok {
			fieldTypeMsgCfg, err := fieldTypeFile.Config.MessageConfig(fieldTypeMessage.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve message %s config", fieldTypeMessage.Name)
			}

			if fieldTypeMsgCfg.UnwrapField {
				unwrappedField, err := g.getUnwrappedField(field)

				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve unwrapped message %s field", fieldTypeMessage.Name)
				}

				fields = append(fields, *unwrappedField)

				continue
			}
		}

		typ, err := g.inputMessageFieldTypeResolver(fieldTypeFile, field)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve `%s.%s` field type", msg.Name, field.Name)
		}

		fields = append(fields, graphql.ObjectField{
			Name:          field.Name,
			Type:          typ,
			QuotedComment: field.QuotedComment,
		})
	}
	for _, field := range msg.MapFields {
		typ, err := g.inputObjectMapFieldTypeResolver(file, field.Map)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve field type")
		}
		fields = append(fields, graphql.ObjectField{
			Name:          field.Name,
			Type:          typ,
			QuotedComment: field.QuotedComment,
		})
	}
	for _, oneOf := range msg.OneOffs {
		for _, fld := range oneOf.Fields {
			fieldTypeFile, err := g.parsedFile(fld.Type.File())
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve value type file")
			}
			typ, err := g.inputMessageFieldTypeResolver(fieldTypeFile, fld)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve field type")
			}
			fields = append(fields, graphql.ObjectField{
				Name:          fld.Name,
				Type:          typ,
				QuotedComment: fld.QuotedComment,
			})
		}
	}

	return fields, nil
}

func (g *Proto2GraphQL) getUnwrappedField(field *parser.Field) (*graphql.ObjectField, error) {
	unwrappedMessage, ok := field.Type.(*parser.Message)

	if !ok {
		return nil, errors.New("can't unwrap field. Not a message")
	}

	unwrappedMessageFile, err := g.parsedFile(unwrappedMessage.File())

	if len(unwrappedMessage.Fields) != 1 {
		return nil, errors.Wrapf(err, "can't unwrap message. Must have 1 field")
	}

	messageFields, err := g.getMessageFields(unwrappedMessageFile, unwrappedMessage)
	if err != nil {
		return nil, err
	}

	unwrappedField := &messageFields[0]
	unwrappedField.Name = field.Name

	if field.Repeated {
		unwrappedField.Type = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(unwrappedField.Type))
	}

	return unwrappedField, nil
}
