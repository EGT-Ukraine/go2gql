package swagger2gql

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/names"
)

type FieldConfig struct {
	ContextKey string `mapstructure:"context_key"`
}
type ObjectConfig struct {
	Fields      map[string]FieldConfig   `mapstructure:"fields"`
	DataLoaders []dataloader.FieldConfig `mapstructure:"data_loaders"`
}

type MethodConfig struct {
	Alias              string         `mapstructure:"alias"`
	RequestType        string         `mapstructure:"request_type"` // QUERY | MUTATION
	DataLoaderProvider ProviderConfig `mapstructure:"data_loader_provider"`
}

type ProviderConfig struct {
	Name         string        `mapstructure:"name"`
	WaitDuration time.Duration `mapstructure:"wait_duration"`
	Slice        bool
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

type ParamConfig struct {
	ParamName  string `mapstructure:"param_name"`
	ContextKey string `mapstructure:"context_key"`
}

type SwaggerFileConfig struct {
	Name string `mapstructure:"name"`

	Path string `mapstructure:"path"`

	ModelsGoPath string `mapstructure:"models_go_path"`

	OutputPkg  string `mapstructure:"output_package"`
	OutputPath string `mapstructure:"output_path"`

	GQLObjectsPrefix string `mapstructure:"gql_objects_prefix"`

	Tags         map[string]*TagConfig     `mapstructure:"tags"`
	Objects      []map[string]ObjectConfig `mapstructure:"objects"`
	ParamsConfig []ParamConfig             `mapstructure:"params_config"`
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

func (pc *SwaggerFileConfig) FieldConfig(objName string, fieldName string) (FieldConfig, error) {
	cfg, err := pc.ObjectConfig(objName)

	if err != nil {
		return FieldConfig{}, errors.Wrap(err, "failed to resolve property config")
	}

	if cfg.Fields != nil {
		paramGqlName := names.FilterNotSupportedFieldNameCharacters(fieldName)

		paramCfg, ok := cfg.Fields[paramGqlName]

		if ok {
			return paramCfg, nil
		}
	}

	for _, paramConfig := range pc.ParamsConfig {
		if paramConfig.ParamName == fieldName {
			return FieldConfig{ContextKey: paramConfig.ContextKey}, nil
		}
	}

	return FieldConfig{}, nil
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
