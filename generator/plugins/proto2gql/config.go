package proto2gql

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
)

const (
	RequestTypeQuery    = "QUERY"
	RequestTypeMutation = "MUTATION"
)

type FieldsConfig struct {
	ContextKey string `mapstructure:"context_key"`
}
type MessageConfig struct {
	ErrorField  string                   `mapstructure:"error_field"`
	Fields      map[string]FieldsConfig  `mapstructure:"fields"`
	DataLoaders []dataloader.FieldConfig `mapstructure:"data_loaders"`
}
type MethodConfig struct {
	Alias               string                    `mapstructure:"alias"`
	RequestType         string                    `mapstructure:"request_type"` // QUERY | MUTATION
	DataLoaderProvider  dataloader.ProviderConfig `mapstructure:"data_loader_provider"`
	UnwrapResponseField bool                      `mapstructure:"unwrap_response_field"`
}
type ServiceConfig struct {
	ServiceName string                  `mapstructure:"service_name"`
	Methods     map[string]MethodConfig `mapstructure:"methods"`
}
type Config struct {
	Files []*ProtoFileConfig `mapstructure:"files"`

	// Global configs for proto files
	Paths          []string                   `mapstructure:"paths"`
	OutputPath     string                     `mapstructure:"output_path"`
	ImportsAliases []map[string]string        `mapstructure:"imports_aliases"`
	Messages       []map[string]MessageConfig `mapstructure:"messages"`
}

func (c *Config) GetOutputPath() string {
	if c == nil {
		return ""
	}

	return c.OutputPath
}

type ProtoFileConfig struct {
	Name string `mapstructure:"name"`

	Paths          []string            `mapstructure:"paths"`
	ImportsAliases []map[string]string `mapstructure:"imports_aliases"`

	ProtoPath string `mapstructure:"proto_path"`

	OutputPkg  string `mapstructure:"output_package"`
	OutputPath string `mapstructure:"output_path"`

	ProtoGoPackage string `mapstructure:"proto_go_package"` // go package of protoc generated code

	GQLEnumsPrefix   string `mapstructure:"gql_enums_prefix"`
	GQLMessagePrefix string `mapstructure:"gql_messages_prefix"`

	Services map[string]ServiceConfig   `mapstructure:"services"`
	Messages []map[string]MessageConfig `mapstructure:"messages"`
}

func (pc *ProtoFileConfig) MessageConfig(msgName string) (MessageConfig, error) {
	if pc == nil {
		return MessageConfig{}, nil
	}
	for _, cfgs := range pc.Messages {
		for msgNameRegex, cfg := range cfgs {
			r, err := regexp.Compile(msgNameRegex)
			if err != nil {
				return MessageConfig{}, errors.Wrapf(err, "failed to compile message name regex '%s'", msgNameRegex)
			}
			if r.MatchString(msgName) {
				return cfg, nil
			}
		}
	}

	return MessageConfig{}, nil
}

func (pc *ProtoFileConfig) GetName() string {
	if pc == nil {
		return ""
	}

	return pc.Name
}

func (pc *ProtoFileConfig) GetPaths() []string {
	if pc == nil {
		return []string{}
	}

	return pc.Paths
}

func (pc *ProtoFileConfig) GetProtoPath() string {
	if pc == nil {
		return ""
	}

	return pc.ProtoPath
}

func (pc *ProtoFileConfig) GetOutputPkg() string {
	if pc == nil {
		return ""
	}

	return pc.OutputPkg
}

func (pc *ProtoFileConfig) GetGoPackage() string {
	if pc == nil {
		return ""
	}

	return pc.ProtoGoPackage
}

func (pc *ProtoFileConfig) GetOutputPath() string {
	if pc == nil {
		return ""
	}

	return pc.OutputPath
}

func (pc *ProtoFileConfig) GetGQLEnumsPrefix() string {
	if pc == nil {
		return ""
	}

	return pc.GQLEnumsPrefix
}

func (pc *ProtoFileConfig) GetGQLMessagePrefix() string {
	if pc == nil {
		return ""
	}

	return pc.GQLMessagePrefix
}

func (pc *ProtoFileConfig) GetImportsAliases() []map[string]string {
	if pc == nil {
		return []map[string]string{}
	}

	return pc.ImportsAliases
}

func (pc *ProtoFileConfig) GetServices() map[string]ServiceConfig {
	if pc == nil {
		return map[string]ServiceConfig{}
	}

	return pc.Services
}

func (pc *ProtoFileConfig) GetMessages() []map[string]MessageConfig {
	if pc == nil {
		return []map[string]MessageConfig{}
	}

	return pc.Messages
}

type SchemaNodeConfig struct {
	Type           string             `mapstructure:"type"` // "OBJECT|SERVICE"
	Proto          string             `mapstructure:"proto"`
	Service        string             `mapstructure:"service"`
	ObjectName     string             `mapstructure:"object_name"`
	Field          string             `mapstructure:"field"`
	Fields         []SchemaNodeConfig `mapstructure:"fields"`
	ExcludeMethods []string           `mapstructure:"exclude_methods"`
	FilterMethods  []string           `mapstructure:"filter_methods"`
}
type SchemaConfig struct {
	Name          string            `mapstructure:"name"`
	OutputPath    string            `mapstructure:"output_path"`
	OutputPackage string            `mapstructure:"output_package"`
	Queries       *SchemaNodeConfig `mapstructure:"queries"`
	Mutations     *SchemaNodeConfig `mapstructure:"mutations"`
}
type GenerateConfig struct {
	Tracer     bool
	VendorPath string
}
