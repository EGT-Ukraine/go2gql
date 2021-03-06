package swagger2gql

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/names"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

var ErrMultipleSuccessResponses = errors.New("method  contains multiple success responses")

func (p *Plugin) graphqlMethod(tagCfg TagConfig, methodCfg MethodConfig, file *parsedFile, tag parser.Tag, method parser.Method) (*graphql.Method, error) {
	name := method.OperationID
	if methodCfg.Alias != "" {
		name = methodCfg.Alias
	}
	var successResponse parser.MethodResponse
	var successResponseFound bool
	for _, resp := range method.Responses {
		if resp.StatusCode/100 == 2 {
			if successResponseFound {
				return nil, ErrMultipleSuccessResponses
			}

			successResponse = resp
			successResponseFound = true
		}
	}
	responseType, err := p.TypeOutputTypeResolver(file, successResponse.ResultType, false)
	if err != nil {
		return nil, errors.Wrap(err, "can't get response output type resolver")
	}
	gqlInputObjName := p.methodParamsInputObjectGQLName(file, method)
	var args []graphql.MethodArgument
	for _, param := range method.Parameters {
		gqlName := names.FilterNotSupportedFieldNameCharacters(param.Name)

		paramCfg, err := file.Config.FieldConfig(gqlInputObjName, param.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve property %s config", param.Name)
		}

		if paramCfg.ContextKey != "" {
			continue
		}
		paramType, err := p.TypeInputTypeResolver(file, param.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve parameter '%s' type resolver", param.Name)
		}

		args = append(args, graphql.MethodArgument{
			Name:          gqlName,
			Type:          paramType,
			QuotedComment: strconv.Quote(param.Description),
		})
	}
	reqType := graphql.GoType{
		Kind: reflect.Ptr,
		ElemType: &graphql.GoType{
			Kind: reflect.Interface,
			Pkg:  file.Config.Tags[tag.Name].ClientGoPackage,
			Name: pascalize(method.OperationID) + "Params",
		},
	}

	if err := p.addDataLoaderProvider(methodCfg, tag, tagCfg, method, successResponse.ResultType, file); err != nil {
		return nil, errors.Wrap(err, "failed add data loader provider")
	}

	return &graphql.Method{
		OriginalName:           method.Path,
		Name:                   name,
		QuotedComment:          strconv.Quote(method.Description),
		GraphQLOutputType:      responseType,
		Arguments:              args,
		RequestResolver:        graphql.ResolverCall(file.OutputPkg, "Resolve"+pascalize(method.OperationID)+"Params"),
		RequestResolverWithErr: true,
		ClientMethodCaller: func(client, req string, ctx graphql.BodyContext) string {
			var res string
			var err error
			if successResponse.ResultType.Kind() == parser.KindNull {
				res, err = p.renderNullMethodCaller(reqType.String(ctx.Importer), req, client, pascalize(method.OperationID))
			} else {
				respType, err := p.goTypeByParserType(file, successResponse.ResultType, true)
				if err != nil {
					panic(errors.Wrap(err, "failed to resolve result go type"))
				}
				res, err = p.renderMethodCaller(respType.String(ctx.Importer), reqType.String(ctx.Importer), req, client, pascalize(method.OperationID))
				if err != nil {
					panic(errors.Wrap(err, "failed to render method caller"))
				}
			}
			if err != nil {
				panic(errors.Wrap(err, "failed to render method caller"))
			}
			return res
		},
		RequestType:          reqType,
		PayloadErrorChecker:  nil,
		PayloadErrorAccessor: nil,
	}, nil
}

func (p *Plugin) addDataLoaderProvider(
	methodCfg MethodConfig,
	tag parser.Tag,
	tagCfg TagConfig,
	method parser.Method,
	successResponseResultType parser.Type,
	file *parsedFile) error {
	dataLoaderProviderConfig := methodCfg.DataLoaderProvider

	if dataLoaderProviderConfig.Name == "" {
		return nil
	}

	resType, ok := successResponseResultType.(*parser.Array)

	if !ok {
		return errors.New("response type must be array")
	}

	responseGoType, err := p.goTypeByParserType(file, resType.ElemType, true)

	if err != nil {
		return err
	}

	dataLoaderOutType, err := p.TypeOutputTypeResolver(file, resType.ElemType, false)

	if err != nil {
		return err
	}

	if len(method.Parameters) != 1 {
		return errors.Errorf("Method %s %s must have 1 input parameter", method.HTTPMethod, method.Path)
	}

	inputArgumentGoType, err := p.goTypeByParserType(file, method.Parameters[0].Type, true)

	if err != nil {
		return err
	}

	if inputArgumentGoType.Kind != reflect.Slice {
		return errors.Errorf("Method %s %s input parameter must be array", method.HTTPMethod, method.Path)
	}

	if !inputArgumentGoType.ElemType.Scalar {
		return errors.Errorf("Method %s %s input parameter must be array of scalars", method.HTTPMethod, method.Path)
	}

	var handleType func(typ parser.Type) string

	handleType = func(typ parser.Type) string {
		switch t := typ.(type) {
		case *parser.Object:
			return p.outputObjectGQLName(file, t)
		case *parser.Array:
			return handleType(t.ElemType)
		}

		panic("Unexpected parser type")
	}

	outputGraphqlTypeName := handleType(resType.ElemType)

	dataLoaderProvider := dataloader.LoaderModel{
		OutputGraphqlTypeName: outputGraphqlTypeName,
		Name:                  dataLoaderProviderConfig.Name,
		WaitDuration:          dataLoaderProviderConfig.WaitDuration,
		Service: &dataloader.Service{
			Name:          p.tagName(tag, &tagCfg),
			CallInterface: p.serviceCallInterface(&tagCfg),
		},
		FetchCode: func(importer *importer.Importer) string {
			elemType := graphql.GoType{
				Kind: reflect.Interface,
				Pkg:  tagCfg.ClientGoPackage,
				Name: pascalize(method.OperationID) + "Params",
			}

			paramsType := elemType.String(importer)
			argName := ucFirst(method.Parameters[0].Name)
			methodName := ucFirst(method.OperationID)

			return `
			params := &` + paramsType + `{
				` + argName + `: keys,
				Context: ctx,
			}

			response, err := client.` + methodName + `(params)

			if err != nil {
				return nil, []error{err}
			}

			return response.Payload, nil
			`
		},
		InputGoType:       inputArgumentGoType,
		OutputGoType:      responseGoType,
		OutputGraphqlType: dataLoaderOutType,
		Slice:             dataLoaderProviderConfig.Slice,
	}

	p.dataLoaderPlugin.AddLoader(dataLoaderProvider)

	return nil
}

func ucFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func (p *Plugin) tagQueriesMethods(tagCfg TagConfig, file *parsedFile, tag parser.Tag) ([]graphql.Method, error) {
	var res []graphql.Method
	for _, method := range tag.Methods {
		var methodCfg MethodConfig
		if tagCfg.Methods[method.Path] != nil {
			methodCfg = tagCfg.Methods[method.Path][strings.ToLower(method.HTTPMethod)]
		}
		if methodCfg.RequestType == "" {
			if method.HTTPMethod != "GET" {
				continue
			}
		} else if methodCfg.RequestType != "QUERY" {
			continue
		}

		meth, err := p.graphqlMethod(tagCfg, methodCfg, file, tag, method)
		if err != nil {
			if err == ErrMultipleSuccessResponses {
				fmt.Println("Warning: Method: ", method.Path, "have multiple successful responses. I'll skip it")
				continue
			}
			return nil, errors.Wrap(err, "failed to resolve graphql method")
		}
		res = append(res, *meth)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res, nil
}
func (p *Plugin) tagMutationsMethods(tagCfg TagConfig, file *parsedFile, tag parser.Tag) ([]graphql.Method, error) {
	var res []graphql.Method
	for _, method := range tag.Methods {
		var methodCfg MethodConfig
		if tagCfg.Methods[method.Path] != nil {
			methodCfg = tagCfg.Methods[method.Path][strings.ToLower(method.HTTPMethod)]
		}
		if methodCfg.RequestType == "" {
			if method.HTTPMethod == "GET" {
				continue
			}
		} else if methodCfg.RequestType != "MUTATION" {
			continue
		}
		meth, err := p.graphqlMethod(tagCfg, methodCfg, file, tag, method)
		if err != nil {
			if err == ErrMultipleSuccessResponses {
				fmt.Println("Warning: Method: ", method.Path, "have multiple successful responses. I'll skip it")
				continue
			}
			return nil, errors.Wrap(err, "failed to resolve graphql method")
		}
		res = append(res, *meth)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res, nil
}
func (p *Plugin) fileServices(file *parsedFile) ([]graphql.Service, error) {
	var res []graphql.Service
	for _, tag := range file.File.Tags {
		tagCfg, ok := file.Config.Tags[tag.Name]

		if !ok {
			fmt.Println("Warning: Skip tag:", tag.Name, "from", file.Config.Path)
			continue
		}

		if tagCfg.ClientGoPackage == "" {
			return nil, errors.Errorf("file: `%s`. Need to specify tag %s `client_go_package` option", file.Config.Name, tag.Name)
		}
		queriesMethods, err := p.tagQueriesMethods(*tagCfg, file, tag)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get tag queries methods")
		}

		mutationsMethods, err := p.tagMutationsMethods(*tagCfg, file, tag)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get tag mutations methods")
		}

		res = append(res, graphql.Service{
			OriginalName:    tag.Name,
			Name:            p.tagName(tag, tagCfg),
			QuotedComment:   strconv.Quote(tag.Description),
			QueryMethods:    queriesMethods,
			MutationMethods: mutationsMethods,
			CallInterface:   p.serviceCallInterface(tagCfg),
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})
	return res, nil
}

func (p *Plugin) tagName(tag parser.Tag, tagCfg *TagConfig) string {
	name := pascalize(tag.Name)

	if tagCfg.ServiceName != "" {
		name = tagCfg.ServiceName
	}

	return name
}

func (p *Plugin) serviceCallInterface(tagCfg *TagConfig) graphql.GoType {
	return graphql.GoType{
		Kind: reflect.Interface,
		Pkg:  tagCfg.ClientGoPackage,
		Name: "IClient",
	}
}
