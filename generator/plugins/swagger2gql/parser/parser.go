package parser

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type handledRefs map[string]Type
type Parser struct {
	parsedFiles []*File
}

func (p Parser) ParsedFiles() []*File {
	return p.parsedFiles
}

func (p *Parser) Parse(loc string, r io.Reader) (*File, error) {
	fullSwaggerFile, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}
	schema := new(spec.Swagger)
	err = schema.UnmarshalJSON(fullSwaggerFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal File")
	}

	fp := fileParser{
		schema:      schema,
		handledRefs: make(map[string]Type),
		result: File{
			file:     schema,
			BasePath: schema.BasePath,
			Location: loc,
		},
	}
	err = fp.parseTags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse file tags")
	}
	return &fp.result, nil
}

type fileParser struct {
	schema      *spec.Swagger
	handledRefs handledRefs
	result      File
}

func resolveScalarType(typ string, format string, enum []interface{}) (Type, error) {
	switch typ {
	case "number":
		switch format {
		case "float":
			return scalarFloat32, nil
		default:
			return scalarFloat64, nil
		}

	case "integer":
		switch format {
		case "int32":
			return scalarInt32, nil
		default:
			return scalarInt64, nil
		}
	case "boolean":
		return scalarBoolean, nil
	case "string":
		if format == "date-time" {
			return ObjDateTime, nil
		}
		if len(enum) > 0 {
			var values = make([]string, len(enum))
			for i, enum := range enum {
				values[i] = enum.(string)
			}
			return scalarString, nil
		}

		return scalarString, nil
	case "file":
		return scalarFile, nil
	}
	return nil, errors.Errorf("scalar type %s is not implemented", typ)
}
func (p *fileParser) resolveSchemaType(route []string, schema *spec.Schema) (Type, error) {
	if schema == nil {
		return &Scalar{kind: KindNull}, nil
	}
	if p.handledRefs == nil {
		p.handledRefs = make(map[string]Type)
	}
	schemaRef := schema.Ref.String()
	if schemaRef != "" {
		if handledType, ok := p.handledRefs[schemaRef]; ok {
			return handledType, nil
		}
		var err error
		schema, err = spec.ResolveRef(p.schema, &schema.Ref)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve $ref")
		}
	}
	if len(schema.Type) != 1 {
		return nil, errors.Errorf("schema type doesn't contains exactly one element: %v", schema.Type)
	}
	switch schema.Type[0] {
	case "array":
		itemSchema := schema.Items.Schema
		itemType, err := p.resolveSchemaType(route, itemSchema)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve array items types")
		}
		return &Array{
			ElemType: itemType,
		}, nil

	case "object":
		if schema.Title != "" {
			route = []string{schema.Title}
		}
		if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
			res := &Map{
				Route: route,
			}
			if schemaRef != "" {
				p.handledRefs[schemaRef] = res
			}
			elemType, err := p.resolveSchemaType(route, schema.AdditionalProperties.Schema)
			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve hashmap value type")
			}
			res.ElemType = elemType

			return res, nil
		}
		typ := &Object{
			Route: route,
			Name:  schema.Title,
		}
		if schemaRef != "" {
			p.handledRefs[schemaRef] = typ
		}
		requiredFields := map[string]struct{}{}
		for _, requiredField := range schema.Required {
			requiredFields[requiredField] = struct{}{}
		}
		for name, prop := range schema.Properties {
			_, required := requiredFields[name]
			ptyp, err := p.resolveSchemaType(append(route, name), &prop)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve prop '%s' type", name)
			}
			typ.Properties = append(typ.Properties, ObjectProperty{
				Name:        name,
				Description: prop.Description,
				Required:    required,
				Type:        ptyp,
			})
		}
		return typ, nil
	}
	return resolveScalarType(schema.Type[0], schema.Format, schema.Enum)
}
func (p *fileParser) parameterType(method *spec.Operation, parameter spec.Parameter) (Type, error) {
	if parameter.Ref.String() != "" || parameter.Schema != nil {
		return p.resolveSchemaType([]string{method.ID}, parameter.Schema)
	}

	if parameter.Type == "array" {
		elemType, err := resolveScalarType(parameter.Items.Type, parameter.Items.Format, parameter.Items.Enum)

		if err != nil {
			return nil, err
		}

		return &Array{
			ElemType: elemType,
		}, nil
	}

	return resolveScalarType(parameter.Type, parameter.Format, parameter.Enum)
}
func (p *fileParser) parseMethodParams(method *spec.Operation) ([]MethodParameter, error) {
	var res []MethodParameter
	for _, parameter := range method.Parameters {
		typ, err := p.parameterType(method, parameter)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve %s parameter type", parameter.Name)
		}
		pos, ok := parameterPositions[parameter.In]
		if !ok {
			return nil, errors.Errorf("unknown parameter position '%s'", parameter.In)
		}
		res = append(res, MethodParameter{
			Name:        parameter.Name,
			Description: parameter.Description,
			Required:    parameter.Required,
			Type:        typ,
			Position:    pos,
		})
	}
	return res, nil
}
func (p *fileParser) parseMethodResponses(method *spec.Operation) ([]MethodResponse, error) {
	var res []MethodResponse
	for statusCode, response := range method.Responses.StatusCodeResponses {
		typ, err := p.resolveSchemaType([]string{method.ID}, response.Schema)
		if err != nil {
			return nil, errors.Wrap(err, "failed to resolve schema type")
		}
		res = append(res, MethodResponse{
			StatusCode:  statusCode,
			Description: response.Description,
			ResultType:  typ,
		})
	}
	return res, nil
}
func (p *fileParser) parseTags() error {
	var tagsByName = make(map[string]*Tag)
	for _, tag := range p.schema.Tags {
		tagsByName[tag.Name] = &Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}
	}
	if p.schema.Paths != nil {
		for path, pathItems := range p.schema.Paths.Paths {
			methods := map[string]*spec.Operation{
				"GET":     pathItems.Get,
				"PUT":     pathItems.Put,
				"POST":    pathItems.Post,
				"DELETE":  pathItems.Delete,
				"OPTIONS": pathItems.Options,
				"HEAD":    pathItems.Head,
				"PATCH":   pathItems.Patch,
			}
			for httpMethod, method := range methods {
				if method == nil {
					continue
				}
				methodTags := method.Tags
				if len(method.Tags) == 0 {
					methodTags = []string{"operations"}
				}

				params, err := p.parseMethodParams(method)
				if err != nil {
					return errors.Wrap(err, "failed to parse method params")
				}
				resps, err := p.parseMethodResponses(method)
				if err != nil {
					return errors.Wrap(err, "failed to resolve method responses")
				}
				m := Method{
					OperationID: method.ID,
					HTTPMethod:  strings.ToUpper(httpMethod),
					Description: method.Description,
					Path:        path,
					Responses:   resps,
					Parameters:  params,
				}
				for _, tag := range methodTags {
					t, ok := tagsByName[tag]
					if !ok {
						t = &Tag{
							Name: tag,
						}
						tagsByName[tag] = t
					}
					t.Methods = append(t.Methods, m)
				}
			}
		}
	}
	var res []Tag
	for _, tag := range tagsByName {
		res = append(res, *tag)
	}
	p.result.Tags = res
	return nil
}
