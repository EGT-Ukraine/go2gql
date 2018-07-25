package graphql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/importer"
)

const (
	PluginName        = "graphql"
	SchemasConfigsKey = "graphql_schemas"
)

type Plugin struct {
	files         map[string]*TypesFile
	schemaConfigs []SchemaConfig
	generateCfg   *generator.GenerateConfig
}

func (p *Plugin) Prepare() error {
	return nil
}
func (p *Plugin) Init(config *generator.GenerateConfig, plugins []generator.Plugin) error {
	var cfgs []SchemaConfig
	p.files = make(map[string]*TypesFile)
	err := mapstructure.Decode(config.PluginsConfigs[SchemasConfigsKey], &cfgs)
	if err != nil {
		return errors.Wrap(err, "failed to decode config")
	}
	p.schemaConfigs = cfgs
	p.generateCfg = config
	return nil
}

// Types returns info about all parsed types
func (p *Plugin) Types() map[string]*TypesFile {
	return p.files
}
func (p *Plugin) AddTypesFile(outputPath string, file *TypesFile) {
	p.files[outputPath] = file
}
func (p Plugin) Name() string {
	return PluginName
}
func (p *Plugin) PrintInfo(infos generator.Infos) {
	if infos.Contains("gql-services") {
		for path, file := range p.files {
			if len(file.Services) > 0 {
				fmt.Println(path)
				for _, service := range file.Services {
					fmt.Println("\t Service " + service.Name)
				}
			}
		}
	}
}
func (p *Plugin) Infos() map[string]string {
	return map[string]string{
		"gql-services": "Info about all parsed GraphQL services",
	}
}
func (p *Plugin) generateTypes() error {
	for outputPath, file := range p.files {
		err := os.MkdirAll(filepath.Dir(outputPath), 0777)
		if err != nil {
			return errors.Wrapf(err, "failed to create directories for output types file %s", outputPath)
		}
		out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s for write", outputPath)
		}
		err = typesGenerator{
			File:          file,
			tracerEnabled: p.generateCfg.GenerateTraces,
			imports: &importer.Importer{
				CurrentPackage: file.Package,
			},
		}.generate(out)
		if err != nil {
			if cerr := out.Close(); cerr != nil {
				err = errors.Wrap(err, cerr.Error())
			}
			return errors.Wrapf(err, "failed to generate types file %s", outputPath)
		}
		if err = out.Close(); err != nil {
			return errors.Wrapf(err, "failed to close generated types file %s", outputPath)
		}
	}
	return nil
}

func (p *Plugin) Generate() error {
	err := p.generateTypes()
	if err != nil {
		return errors.Wrap(err, "failed to generate types files")
	}
	err = p.generateSchemas()
	if err != nil {
		return errors.Wrap(err, "failed to generate schema files")
	}
	return nil
}
