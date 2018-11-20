package proto2gql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql/parser"
)

func (g Proto2GraphQL) serviceMethodArguments(file *parsedFile, cfg MethodConfig, method *parser.Method) ([]graphql.MethodArgument, error) {
	var args []graphql.MethodArgument
	for _, field := range method.InputMessage.Fields {
		typeFile, err := g.parsedFile(field.Type.File())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve field '%s' file", field.Name)
		}
		msgCfg, err := file.Config.MessageConfig(dotedTypeName(method.InputMessage.TypeName))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve message %s config", method.InputMessage.TypeName)
		}
		fldCfg := msgCfg.Fields[field.Name]
		if fldCfg.ContextKey != "" {
			continue
		}
		typResolver, err := g.TypeInputTypeResolver(typeFile, field.Type)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare input type resolver")
		}
		if field.Repeated {
			typResolver = graphql.GqlListTypeResolver(graphql.GqlNonNullTypeResolver(typResolver))
		}
		args = append(args, graphql.MethodArgument{
			Name: field.Name,
			Type: typResolver,
		})
	}
	for _, field := range method.InputMessage.MapFields {
		typeFile, err := g.parsedFile(field.Map.File())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve field '%s' file", field.Name)
		}
		typResolver, err := g.TypeInputTypeResolver(typeFile, field.Map)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare input type resolver")
		}
		args = append(args, graphql.MethodArgument{
			Name: field.Name,
			Type: typResolver,
		})
	}
	return args, nil
}
func (g Proto2GraphQL) messagePayloadErrorParams(message *parser.Message) (checker graphql.PayloadErrorChecker, accessor graphql.PayloadErrorAccessor, err error) {
	outMsgCfg, err := g.fileConfig(message.File()).MessageConfig(dotedTypeName(message.TypeName))
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
	for _, fld := range message.Fields {
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
func (g Proto2GraphQL) serviceMethod(cfg MethodConfig, file *parsedFile, method *parser.Method) (*graphql.Method, error) {
	outputMsgTypeFile, err := g.parsedFile(method.OutputMessage.File())
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve file type file")
	}
	outType, err := g.TypeOutputTypeResolver(outputMsgTypeFile, method.OutputMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get output type resolver for method: %s", method.Name)
	}
	requestType, err := g.goTypeByParserType(method.InputMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get request go type for method: %s", method.Name)
	}
	args, err := g.serviceMethodArguments(file, cfg, method)
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
	valueResolver, valueResolverWithErr, _, err := g.TypeValueResolver(inputMessageFile, method.InputMessage, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve message value resolver")
	}
	return &graphql.Method{
		Name:              g.methodName(cfg, method),
		GraphQLOutputType: outType,
		RequestType:       requestType,
		ClientMethodCaller: func(client, arg string, ctx graphql.BodyContext) string {
			return client + "." + camelCase(method.Name) + "(ctx," + arg + ")"
		},
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
		met, err := g.serviceMethod(mc, file, method)
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
		met, err := g.serviceMethod(mc, file, method)
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
			Name: g.serviceName(sc, service),
			CallInterface: graphql.GoType{
				Kind: reflect.Interface,
				Pkg:  file.GRPCSourcesPkg,
				Name: service.Name + "Client",
			},
			QueryMethods:    queryMethods,
			MutationMethods: mutationsMethods,
		})
	}
	return res, nil
}
