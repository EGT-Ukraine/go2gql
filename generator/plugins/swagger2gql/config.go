package swagger2gql

import (
	"regexp"

	"github.com/pkg/errors"
)

type FieldConfig struct {
	ContextKey string `mapstructure:"context_key"`
}
type ObjectConfig struct {
	ErrorField string                 `mapstructure:"error_field"`
	Fields     map[string]FieldConfig `mapstructure:"fields"`
}

type MethodConfig struct {
	Alias       string `mapstructure:"alias"`
	RequestType string `mapstructure:"request_type"` // QUERY | MUTATION
}
type TagConfig struct {
	ClientGoPackage string                             `mapstructure:"client_go_package"`
	ServiceName     string                             `mapstructure:"service_name"`
	Methods         map[string]map[string]MethodConfig `mapstructure:"methods"`
}
type Config struct {
	Files      []*SwaggerFileConfig      `mapstructure:"files"`
	OutputPath string                    `mapstructure:"output_path"`
	Messages   []map[string]ObjectConfig `mapstructure:"messages"`
}

func (c *Config) GetOutputPath() string {
	if c == nil {
		return ""
	}

	return c.OutputPath
}

type SwaggerFileConfig struct {
	Name string `mapstructure:"name"`

	Path string `mapstructure:"path"`

	ModelsGoPath string `mapstructure:"models_go_path"`

	OutputPkg  string `mapstructure:"output_package"`
	OutputPath string `mapstructure:"output_path"`

	GQLObjectsPrefix string `mapstructure:"gql_objects_prefix"`

	Tags    map[string]*TagConfig     `mapstructure:"tags"`
	Objects []map[string]ObjectConfig `mapstructure:"objects"`
}

func (pc *SwaggerFileConfig) ObjectConfig(objName string) (ObjectConfig, error) {
	if pc == nil {
		return ObjectConfig{}, nil
	}
	for _, cfgs := range pc.Objects {
		for msgNameRegex, cfg := range cfgs {
			r, err := regexp.Compile(msgNameRegex)
			if err != nil {
				return ObjectConfig{}, errors.Wrapf(err, "failed to compile object name regex '%s'", msgNameRegex)
			}
			if r.MatchString(objName) {
				return cfg, nil
			}
		}
	}

	return ObjectConfig{}, nil
}

func (pc *SwaggerFileConfig) GetName() string {
	if pc == nil {
		return ""
	}

	return pc.Name
}

func (pc *SwaggerFileConfig) GetPath() string {
	if pc == nil {
		return ""
	}

	return pc.Path
}

func (pc *SwaggerFileConfig) GetOutputPkg() string {
	if pc == nil {
		return ""
	}

	return pc.OutputPkg
}

func (pc *SwaggerFileConfig) GetOutputPath() string {
	if pc == nil {
		return ""
	}

	return pc.OutputPath
}

func (pc *SwaggerFileConfig) GetGQLMessagePrefix() string {
	if pc == nil {
		return ""
	}

	return pc.GQLObjectsPrefix
}

func (pc *SwaggerFileConfig) GetTags() map[string]*TagConfig {
	if pc == nil {
		return map[string]*TagConfig{}
	}

	return pc.Tags
}

func (pc *SwaggerFileConfig) GetObjects() []map[string]ObjectConfig {
	if pc == nil {
		return []map[string]ObjectConfig{}
	}

	return pc.Objects
}
