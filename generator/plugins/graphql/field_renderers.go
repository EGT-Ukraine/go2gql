package graphql

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type fieldsRenderer struct {
	templateFuncs map[string]interface{}
}

type mapFieldsRenderer struct {
	templateFuncs map[string]interface{}
}

type RenderFieldsContext struct {
	OutputObject  OutputObject
	ObjectContext BodyContext
}

func (r *fieldsRenderer) RenderFields(o OutputObject, ctx BodyContext) (string, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesOutput_fieldsGohtmlBytes()
	if err != nil {
		return "", errors.Wrap(err, "failed to get output fields template")
	}
	bodyTpl, err := template.New("fieldsRenderer").Funcs(r.templateFuncs).Parse(string(tmpl))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}
	err = bodyTpl.Execute(buf, RenderFieldsContext{
		OutputObject:  o,
		ObjectContext: ctx,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return buf.String(), nil
}

func (r *mapFieldsRenderer) RenderFields(o OutputObject, ctx BodyContext) (string, error) {
	buf := new(bytes.Buffer)
	tmpl, err := templatesOutput_map_fieldsGohtmlBytes()
	if err != nil {
		return "", errors.Wrap(err, "failed to get output fields template")
	}
	bodyTpl, err := template.New("mapFieldsRenderer").Funcs(r.templateFuncs).Parse(string(tmpl))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}
	err = bodyTpl.Execute(buf, RenderFieldsContext{
		OutputObject:  o,
		ObjectContext: ctx,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return buf.String(), nil
}
