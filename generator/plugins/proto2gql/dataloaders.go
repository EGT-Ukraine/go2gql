package proto2gql

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g Proto2GraphQL) registerMethodDataLoaders(sc ServiceConfig, cfg MethodConfig, file *parsedFile, method *parser.Method) error {
	for name, cfg := range cfg.DataLoaderProvider {
		if err := g.registerMethodDataLoader(name, cfg, sc, file, method); err != nil {
			return errors.Wrapf(err, "failed to register %s data loader", name)
		}
	}

	return nil
}

func (g Proto2GraphQL) registerMethodDataLoader(name string, cfg DataLoaderConfig, serviceConfig ServiceConfig, file *parsedFile, method *parser.Method) error {
	err := g.validateDataLoader(cfg, method)
	if err != nil {
		return errors.Wrap(err, "validation error")
	}

	resultField, err := messageFieldByPath(method.OutputMessage, cfg.ResultField)
	if err != nil {
		return errors.Wrap(err, "failed to get result field")
	}
	normalResultField := resultField.(*parser.NormalField)
	resultMessage := normalResultField.Type.(*parser.Message)

	matchField, err := messageFieldByPath(resultField.GetType().(*parser.Message), cfg.MatchField)
	if err != nil {
		return errors.Wrap(err, "failed to get match field")
	}

	responseGoType, err := g.goTypeByParserType(normalResultField.Type)
	if err != nil {
		return errors.Wrap(err, "failed to get result field go type")
	}
	outputMsgTypeFile, err := g.parsedFile(method.OutputMessage.File())
	if err != nil {
		return errors.Wrap(err, "failed to resolve file type file")
	}

	dataLoaderOutType, err := g.TypeOutputGraphQLTypeResolver(outputMsgTypeFile, resultMessage)
	if err != nil {
		return errors.Wrap(err, "failed to resolve output type")
	}

	var fetchCode func(importer *importer.Importer) string
	if cfg.Type == DataLoaderType1ToN {
		fetchCode, err = g.oneToNDataLoaderFetchCode(file, cfg, method)
		if err != nil {
			return errors.Wrap(err, "failed to resolve 1-N data loader fetch code")
		}
		responseCopy := responseGoType
		responseGoType = graphql.GoType{
			Kind:     reflect.Slice,
			ElemType: &responseCopy,
		}
		dataLoaderOutType = graphql.GqlListTypeResolver(dataLoaderOutType)
	} else {
		fetchCode, err = g.oneToOneDataLoaderFetchCode(file, cfg, method)
		if err != nil {
			return errors.Wrap(err, "failed to resolve 1-1 data loader fetch code")
		}
	}

	matchFieldGoType, err := g.goTypeByParserType(matchField.GetType())
	if err != nil {
		return errors.Wrap(err, "failed to resolve go type of match field")
	}

	dataLoaderProvider := dataloader.LoaderModel{
		Service: &dataloader.Service{
			Name:          g.serviceName(serviceConfig, method.Service),
			CallInterface: g.serviceCallInterface(file, method.Service.Name),
		},
		FetchCode: fetchCode,
		InputGoType: graphql.GoType{
			Kind:     reflect.Slice,
			ElemType: &matchFieldGoType,
		},
		OutputGoType:      responseGoType,
		OutputGraphqlType: dataLoaderOutType,
		Name:              name,
		WaitDuration:      cfg.WaitDuration,
		Slice:             cfg.Type == DataLoaderType1ToN,
	}

	g.DataLoaderPlugin.AddLoader(dataLoaderProvider)

	return nil
}

func (g Proto2GraphQL) validateDataLoader(cfg DataLoaderConfig, method *parser.Method) error {
	if cfg.RequestField == "" {
		return errors.New("empty request field")
	}
	requestField, err := messageFieldByPath(method.InputMessage, cfg.RequestField)
	if err != nil {
		return errors.Wrap(err, "failed to get message request field")
	}
	normalRequestField, ok := requestField.(*parser.NormalField)
	if !ok {
		return errors.New("request field should not be a map")
	}
	if !normalRequestField.Repeated {
		return errors.New("request field should be repeated")
	}

	if cfg.ResultField == "" {
		return errors.New("empty result field")
	}
	resultField, err := messageFieldByPath(method.OutputMessage, cfg.ResultField)
	if err != nil {
		return errors.Wrap(err, "failed to get message result field")
	}
	normalResultField, ok := resultField.(*parser.NormalField)
	if !ok {
		return errors.New("result field should not be a map")
	}
	if !normalResultField.Repeated {
		return errors.New("result field should be repeated")
	}
	if normalResultField.Type.Kind() != parser.TypeMessage {
		return errors.New("result field should be of message type")
	}
	resultMessage := normalResultField.Type.(*parser.Message)

	if cfg.MatchField == "" {
		return errors.New("empty match field")
	}
	matchField, err := messageFieldByPath(resultMessage, cfg.MatchField)
	if err != nil {
		return errors.Wrap(err, "failed to get message match field")
	}
	normalMatchField, ok := matchField.(*parser.NormalField)
	if !ok {
		return errors.New("match field should not be a map")
	}
	if normalMatchField.Type.Kind() != parser.TypeScalar {
		return errors.New("match field should not be scalar")
	}
	if normalMatchField.Repeated {
		return errors.New("match should not be repeated")
	}

	switch cfg.Type {
	case DataLoaderType1ToN, DataLoaderType1To1:
	default:
		return errors.Errorf("invalid dataloader type(%s)", cfg.Type)
	}

	return nil
}

func (g Proto2GraphQL) oneToOneDataLoaderFetchCode(file *parsedFile, cfg DataLoaderConfig, method *parser.Method) (func(importer *importer.Importer) string, error) {
	resultField, err := messageFieldByPath(method.OutputMessage, cfg.ResultField)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result message field by path")
	}
	responseGoType, err := g.goTypeByParserType(resultField.GetType())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve go type by result field parser type")
	}

	return func(importer *importer.Importer) string {
		filledRequest, err := g.getMessageWithFilledField(file, importer, method.InputMessage, strings.Split(cfg.RequestField, "."), "keys")
		if err != nil {
			panic("failed to build request message")
		}

		responseFieldAccessor := buildFieldAccessorByFieldPath(cfg.ResultField)
		matchFieldAccessor := buildFieldAccessorByFieldPath(cfg.MatchField)

		return `response, err := client.` + method.Name + `(ctx, ` + filledRequest + `)
				if err != nil{
					return nil, []error{err}
				}
				var result = make([]` + responseGoType.String(importer) + `, len(keys))
				for i, key := range keys {
				    for _, value := range response.` + responseFieldAccessor + ` {
				        if value.` + matchFieldAccessor + ` == key {
				            result[i] = value
							break
				        }
				    }
				}
				
				return result, nil`
	}, nil
}

func (g Proto2GraphQL) oneToNDataLoaderFetchCode(file *parsedFile, cfg DataLoaderConfig, method *parser.Method) (func(importer *importer.Importer) string, error) {
	resultField, err := messageFieldByPath(method.OutputMessage, cfg.ResultField)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result message field by path")
	}
	responseGoType, err := g.goTypeByParserType(resultField.GetType())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve go type by result field parser type")
	}

	return func(importer *importer.Importer) string {
		filledRequest, err := g.getMessageWithFilledField(file, importer, method.InputMessage, strings.Split(cfg.RequestField, "."), "keys")
		if err != nil {
			panic("failed to build request message")
		}

		responseFieldAccessor := buildFieldAccessorByFieldPath(cfg.ResultField)
		matchFieldAccessor := buildFieldAccessorByFieldPath(cfg.MatchField)

		return `response, err := client.` + method.Name + `(ctx, ` + filledRequest + `)
				if err != nil{
					return nil, []error{err}
				}
				var result = make([][]` + responseGoType.String(importer) + `, len(keys))
				for i, key := range keys {
				    for _, value := range response.` + responseFieldAccessor + ` {
				        if value.` + matchFieldAccessor + ` == key {
				            result[i] = append(result[i], value)
				        }
				    }
				}
				
				return result, nil`
	}, nil
}

func messageFieldByPath(message *parser.Message, path string) (parser.Field, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("empty path")
	}
	pathParts := strings.Split(path, ".")
	msg := message
	for i := 0; i < len(pathParts); i++ {
		part := pathParts[i]
		var found bool
		for _, field := range msg.GetFields() {
			if field.GetName() == part {
				if i == len(pathParts)-1 {
					return field, nil
				}
				if field.GetType().Kind() != parser.TypeMessage {
					return nil, errors.Errorf("field %s of message %s is not a message", field.GetName(), msg.Name)
				}
				msg = field.GetType().(*parser.Message)
				found = true
			}
		}
		if !found {
			return nil, errors.Errorf("field %s was not found in message %s", part, msg.Name)
		}
	}

	panic("Something went wrong") // We shouldn't reach this line
}

func buildFieldAccessorByFieldPath(path string) string {
	pathParts := strings.Split(path, ".")
	for i, part := range pathParts {
		pathParts[i] = camelCase(part)
	}

	return strings.Join(pathParts, ".")
}

func (g Proto2GraphQL) getMessageWithFilledField(file *parsedFile, importer *importer.Importer, msg *parser.Message, pathParts []string, value string) (string, error) {

	field, ok := msg.GetFieldByName(pathParts[0])
	if !ok {
		return "", errors.New("can't find field ")
	}
	if len(pathParts) == 1 {
		typ, err := g.goTypeByParserType(msg)
		if err != nil {
			return "", errors.Wrap(err, "failed to resolve go type by parser type")
		}

		return g.buildMessageWithFilledField(file, typ.ElemType.String(importer), importer, camelCase(field.GetName()), value), nil
	}
	fieldMsg := field.GetType().(*parser.Message)
	res, err := g.getMessageWithFilledField(file, importer, fieldMsg, pathParts[1:], value)
	if err != nil {
		return "", errors.Wrapf(err, "failed to generate message %s with filled %v", fieldMsg.Name, pathParts[1:])
	}
	typ, err := g.goTypeByParserType(msg)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve go type by parser type")
	}

	return g.buildMessageWithFilledField(file, typ.ElemType.String(importer), importer, camelCase(field.GetName()), res), nil
}

func (g Proto2GraphQL) buildMessageWithFilledField(file *parsedFile, typ string, importer *importer.Importer, field, value string) string {
	return `&` + typ + `{` + field + `: ` + value + `}`
}
