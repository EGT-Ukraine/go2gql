package swagger2gql

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
)

func (p *Plugin) renderArrayValueResolver(arg string, resultGoTyp graphql.GoType, ctx graphql.BodyContext, elemResolver graphql.ValueResolver, elemResolverWithErr bool) (string, error) {
	tplBody, err := templatesValue_resolver_arrayGohtmlBytes()
	if err != nil {
		panic(errors.Wrap(err, "failed to get array value resolver template").Error())
	}
	tpl, err := template.New("array_value_resolver").Funcs(template.FuncMap{
		"errorsPkg": func() string {
			return ctx.Importer.New(graphql.ErrorsPkgPath)
		},
	}).Parse(string(tplBody))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse array value resolver template"))
	}
	res := new(bytes.Buffer)
	err = tpl.Execute(res, map[string]interface{}{
		"resultType": func() string {
			return resultGoTyp.String(ctx.Importer)
		},
		"rootCtx":             ctx,
		"elemResolver":        elemResolver,
		"elemResolverWithErr": elemResolverWithErr,
		"arg": arg,
	})
	return res.String(), err
}

func (p *Plugin) renderPtrDatetimeResolver(arg string, ctx graphql.BodyContext) (string, error) {
	tplBody, err := templatesValue_resolver_ptr_datetimeGohtmlBytes()
	if err != nil {
		panic(errors.Wrap(err, "failed to get array value resolver template").Error())
	}
	tpl, err := template.New("array_value_resolver").Funcs(template.FuncMap{
		"strfmtPkg": func() string {
			return ctx.Importer.New(strFmtPkg)
		},
		"errorsPkg": func() string {
			return ctx.Importer.New(graphql.ErrorsPkgPath)
		},
		"timePkg": func() string {
			return ctx.Importer.New(timePkg)
		},
	}).Parse(string(tplBody))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse array value resolver template"))
	}
	res := new(bytes.Buffer)
	err = tpl.Execute(res, map[string]interface{}{
		"arg": arg,
	})
	return res.String(), err
}

func (p *Plugin) renderDatetimeValueResolverTemplate(arg string, ctx graphql.BodyContext) (string, error) {
	tplBody, err := templatesValue_resolver_datetimeGohtmlBytes()
	if err != nil {
		panic(errors.Wrap(err, "failed to get array value resolver template").Error())
	}
	tpl, err := template.New("array_value_resolver").Funcs(template.FuncMap{
		"strfmtPkg": func() string {
			return ctx.Importer.New(strFmtPkg)
		},
		"errorsPkg": func() string {
			return ctx.Importer.New(graphql.ErrorsPkgPath)
		},
		"timePkg": func() string {
			return ctx.Importer.New(timePkg)
		},
	}).Parse(string(tplBody))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse array value resolver template"))
	}
	res := new(bytes.Buffer)
	err = tpl.Execute(res, map[string]interface{}{
		"arg": arg,
	})
	return res.String(), err
}

func (p *Plugin) renderMethodCaller(responseType, requestType, requestVar, clientVar, methodName string) (string, error) {
	tplBody, err := templatesMethod_callerGohtmlBytes()
	if err != nil {
		panic(errors.Wrap(err, "failed to get array value resolver template").Error())
	}
	tpl, err := template.New("array_value_resolver").Parse(string(tplBody))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse array value resolver template"))
	}
	res := new(bytes.Buffer)
	err = tpl.Execute(res, map[string]interface{}{
		"clientVar":  clientVar,
		"methodName": methodName,
		"reqType":    requestType,
		"reqVar":     requestVar,
		"respType":   responseType,
	})
	return res.String(), err
}

func (p *Plugin) renderNullMethodCaller(requestType, requestVar, clientVar, methodName string) (string, error) {
	tplBody, err := templatesMethod_caller_nullGohtmlBytes()
	if err != nil {
		panic(errors.Wrap(err, "failed to get array value resolver template").Error())
	}
	tpl, err := template.New("array_value_resolver").Parse(string(tplBody))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse array value resolver template"))
	}
	res := new(bytes.Buffer)
	err = tpl.Execute(res, map[string]interface{}{
		"clientVar":  clientVar,
		"methodName": methodName,
		"reqType":    requestType,
		"reqVar":     requestVar,
	})
	return res.String(), err
}
