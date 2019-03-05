package proto2gql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g Proto2GraphQL) serviceMethodArguments(file *parsedFile, method *parser.Method) ([]graphql.MethodArgument, error) {
	var args []graphql.MethodArgument

	messageFields, err := g.getMessageFields(file, method.InputMessage)
	if err != nil {
		return nil, err
	}

	for _, messageField := range messageFields {
		args = append(args, graphql.MethodArgument{
			Name:          messageField.Name,
			Type:          messageField.Type,
			QuotedComment: messageField.QuotedComment,
		})
	}

	return args, nil
}
func (g Proto2GraphQL) messagePayloadErrorParams(message *parser.Message) (checker graphql.PayloadErrorChecker, accessor graphql.PayloadErrorAccessor, err error) {
	outMsgCfg, err := g.fileConfig(message.File()).MessageConfig(message.Name)
	if err != nil {
		err = errors.Wrap(err, "failed to resolve output message config")

		return
	}
	if outMsgCfg.ErrorField == "" {
		return
	}
	errorAccessor := func(arg string) string {
		return arg + ".Get" + camelCase(outMsgCfg.ErrorField) + "()"
	}
	errorCheckerByType := func(repeated bool, p parser.Type) graphql.PayloadErrorChecker {
		if repeated || p.Kind() == parser.TypeMap {
			return func(arg string) string {
				return "len(" + arg + ".Get" + camelCase(outMsgCfg.ErrorField) + "())>0"
			}
		}
		if p.Kind() == parser.TypeScalar || p.Kind() == parser.TypeEnum {
			fmt.Println("Warning: scalars and enums is not supported as payload error fields")

			return nil
		}
		if p.Kind() == parser.TypeMessage {
			return func(arg string) string {
				return arg + ".Get" + camelCase(outMsgCfg.ErrorField) + "() != nil"
			}
		}

		return nil
	}
	for _, fld := range message.NormalFields {
		if fld.Name == outMsgCfg.ErrorField {
			errorChecker := errorCheckerByType(fld.Repeated, fld.Type)
			if errorChecker == nil {
				return nil, nil, nil
			}

			return errorChecker, errorAccessor, nil
		}
	}
	for _, fld := range message.MapFields {
		if fld.Name == outMsgCfg.ErrorField {
			errorChecker := errorCheckerByType(false, fld.Map)
			if errorChecker == nil {
				return nil, nil, nil
			}

			return errorChecker, errorAccessor, nil
		}
	}
	for _, of := range message.OneOffs {
		for _, fld := range of.Fields {
			if fld.Name == outMsgCfg.ErrorField {
				errorChecker := errorCheckerByType(false, fld.Type)
				if errorChecker == nil {
					return nil, nil, nil
				}

				return errorChecker, errorAccessor, nil
			}
		}
	}

	return nil, nil, nil
}

func (g Proto2GraphQL) methodName(cfg MethodConfig, method *parser.Method) string {
	if cfg.Alias != "" {
		return cfg.Alias
	}

	return method.Name
}

func (g Proto2GraphQL) serviceMethod(sc ServiceConfig, cfg MethodConfig, file *parsedFile, method *parser.Method) (*graphql.Method, error) {
	outputMsgTypeFile, err := g.parsedFile(method.OutputMessage.File())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve file type file")
	}

	clientMethodCaller := func(client, arg string, ctx graphql.BodyContext) string {
		return client + "." + camelCase(method.Name) + "(ctx," + arg + ")"
	}

	var outProtoType parser.Type
	var outType graphql.TypeResolver

	var outProtoTypeRepeated bool

	outputMessageConfig, err := file.Config.MessageConfig(method.OutputMessage.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve message %s config", method.OutputMessage.Name)
	}

	if outputMessageConfig.UnwrapField {
		if len(method.OutputMessage.NormalFields) != 1 {
			return nil, errors.Errorf(
				"can't unwrap `%s` service `%s` method response. Output message must have 1 field.",
				method.Service.Name,
				method.Name,
			)
		}

		unwrapFieldName := method.OutputMessage.NormalFields[0].Name

		resolver, err := g.FieldOutputValueResolver(method.OutputMessage, unwrapFieldName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build output value resolver")
		}

		outType, err = g.FieldOutputGraphQLTypeResolver(method.OutputMessage, unwrapFieldName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build output type resovler")
		}

		clientMethodCaller = func(client, arg string, ctx graphql.BodyContext) string {
			return `func() (interface{}, error) {
				res, err :=  ` + client + "." + camelCase(method.Name) + `(ctx,` + arg + `)

				if err != nil {
					return nil, err
				}

				return ` + resolver("res", ctx) + `, nil
			}()`
		}
	} else {
		if len(method.OutputMessage.NormalFields) == 1 {
			fmt.Printf(
				"Suggestion: service `%s` method `%s` in file `%s` has 1 output field. Can be unwrapped.\n",
				method.Service.Name,
				method.Name,
				file.File.FilePath,
			)
		}

		outProtoType = method.OutputMessage
		outType, err = g.TypeOutputGraphQLTypeResolver(outputMsgTypeFile, outProtoType)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get output type resolver for method: %s", method.Name)
		}
		if outProtoTypeRepeated {
			outType = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(outType))
		}
	}

	requestType, err := g.goTypeByParserType(method.InputMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get request go type for method: %s", method.Name)
	}
	args, err := g.serviceMethodArguments(file, method)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare service method arguments")
	}
	payloadErrChecker, payloadErrAccessor, err := g.messagePayloadErrorParams(method.OutputMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve message payload error params")
	}
	inputMessageFile, err := g.parsedFile(method.InputMessage.File())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve message '%s' parsed file", dotedTypeName(method.InputMessage.TypeName))
	}
	valueResolver, valueResolverWithErr, _, err := g.TypeValueResolver(inputMessageFile, method.InputMessage, "", false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve message value resolver")
	}

	if err := g.registerMethodDataLoaders(sc, cfg, file, method); err != nil {
		return nil, errors.Wrap(err, "failed add data loader provider")
	}

	return &graphql.Method{
		OriginalName:           method.Name,
		Name:                   g.methodName(cfg, method),
		QuotedComment:          method.QuotedComment,
		GraphQLOutputType:      outType,
		RequestType:            requestType,
		ClientMethodCaller:     clientMethodCaller,
		RequestResolver:        valueResolver,
		RequestResolverWithErr: valueResolverWithErr,
		Arguments:              args,
		PayloadErrorChecker:    payloadErrChecker,
		PayloadErrorAccessor:   payloadErrAccessor,
	}, nil
}

func (g Proto2GraphQL) serviceQueryMethods(sc ServiceConfig, file *parsedFile, service *parser.Service) ([]graphql.Method, error) {
	var res []graphql.Method
	for _, method := range service.Methods {
		mc := sc.Methods[method.Name]
		if !g.methodIsQuery(mc, method) {
			continue
		}
		met, err := g.serviceMethod(sc, mc, file, method)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare service method %s", method.Name)
		}
		res = append(res, *met)
	}

	return res, nil
}
func (g Proto2GraphQL) methodIsQuery(cfg MethodConfig, method *parser.Method) bool {
	switch cfg.RequestType {
	case RequestTypeQuery:
		return true
	case RequestTypeMutation:
		return false
	}

	return strings.HasPrefix(strings.ToLower(method.Name), "get")
}

func (g Proto2GraphQL) serviceMutationsMethods(cfg ServiceConfig, file *parsedFile, service *parser.Service) ([]graphql.Method, error) {
	var res []graphql.Method
	for _, method := range service.Methods {
		mc := cfg.Methods[method.Name]
		if g.methodIsQuery(mc, method) {
			continue
		}
		met, err := g.serviceMethod(cfg, mc, file, method)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to prepare service method %s", method.Name)
		}

		res = append(res, *met)
	}

	return res, nil
}
func (g Proto2GraphQL) serviceName(sc ServiceConfig, service *parser.Service) string {
	if sc.ServiceName != "" {
		return sc.ServiceName
	}

	return service.Name
}

func (g Proto2GraphQL) fileServices(file *parsedFile) ([]graphql.Service, error) {
	var res []graphql.Service
	for _, service := range file.File.Services {
		sc := file.Config.GetServices()[service.Name]
		queryMethods, err := g.serviceQueryMethods(sc, file, service)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve service methods")
		}
		mutationsMethods, err := g.serviceMutationsMethods(sc, file, service)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve service methods")
		}

		res = append(res, graphql.Service{
			OriginalName:    service.Name,
			Name:            g.serviceName(sc, service),
			QuotedComment:   service.QuotedComment,
			CallInterface:   g.serviceCallInterface(file, service.Name),
			QueryMethods:    queryMethods,
			MutationMethods: mutationsMethods,
		})
	}

	return res, nil
}

func (g Proto2GraphQL) serviceCallInterface(file *parsedFile, serviceName string) graphql.GoType {
	return graphql.GoType{
		Kind: reflect.Interface,
		Pkg:  file.GRPCSourcesPkg,
		Name: serviceName + "Client",
	}
}
