package swagger2gql

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql/parser"
)

const (
	PluginName            = "swagger2gql"
	PluginConfigKey       = "swagger2gql"
	PluginImportConfigKey = "swagger2gql_files"
)

type Plugin struct {
	graphql          *graphql.Plugin
	dataLoaderPlugin *dataloader.Plugin
	config           *Config
	generateConfig   *generator.GenerateConfig

	parsedFiles []*parsedFile
}

func (p *Plugin) Init(config *generator.GenerateConfig, plugins []generator.Plugin) error {
	for _, plugin := range plugins {
		switch plugin.Name() {
		case graphql.PluginName:
			p.graphql = plugin.(*graphql.Plugin)
		case dataloader.PluginName:
			p.dataLoaderPlugin = plugin.(*dataloader.Plugin)
		}
	}
	if p.graphql == nil {
		return errors.New("'graphql' plugin is not installed.")
	}
	if p.dataLoaderPlugin == nil {
		return errors.New("'dataloader' plugin is not installed.")
	}
	cfg := new(Config)
	err := mapstructure.Decode(config.PluginsConfigs[PluginConfigKey], cfg)
	if err != nil {
		return errors.Wrap(err, "failed to decode config")
	}
	p.generateConfig = config
	p.config = cfg

	if err = p.parseImports(); err != nil {
		return errors.Wrap(err, "failed to decode imports")
	}

	return nil
}

func (p *Plugin) parseImports() error {
	for _, pluginsConfigsImports := range p.generateConfig.PluginsConfigsImports {
		configs := new([]*SwaggerFileConfig)
		if err := mapstructure.Decode(pluginsConfigsImports.PluginsConfigs[PluginImportConfigKey], configs); err != nil {
			return errors.Wrap(err, "failed to decode config")
		}

		if err := p.processConfigs(pluginsConfigsImports.Path, *configs); err != nil {
			return errors.Wrap(err, "failed to process configs")
		}
	}

	return nil
}

func (p *Plugin) processConfigs(path string, configs []*SwaggerFileConfig) error {
	for _, config := range configs {
		var importFileDir = filepath.Dir(path)
		var swaggerPath = filepath.Join(importFileDir, config.Path)

		if config.ModelsGoPath == "" {
			goPkg, err := GoPackageByPath(importFileDir, p.generateConfig.VendorPath)

			if err != nil {
				return errors.Wrap(err, "failed to get go package by path "+importFileDir)
			}

			config.ModelsGoPath = goPkg + "/models"
		}

		if err := p.processTags(config, importFileDir); err != nil {
			return errors.Wrap(err, "failed to process tags")
		}

		config.Path = swaggerPath

		p.config.Files = append(p.config.Files, config)
	}

	return nil
}

func (p *Plugin) processTags(config *SwaggerFileConfig, importFileDir string) error {
	for tagName, tag := range config.Tags {
		goPkg, err := GoPackageByPath(importFileDir, p.generateConfig.VendorPath)

		if err != nil {
			return errors.Wrap(err, "failed to get go package by path "+importFileDir)
		}

		if tag == nil {
			tag = new(TagConfig)
			config.Tags[tagName] = tag
		}

		if tag.ClientGoPackage == "" {
			controllerName := strings.Replace(tagName, "-", "_", -1)

			tag.ClientGoPackage = goPkg + "/client/" + controllerName
		}
	}

	return nil
}

func (p *Plugin) PrintInfo(info generator.Infos) {
}

func (p *Plugin) Infos() map[string]string {
	return nil
}

func (p *Plugin) prepareTypesFile(file *parsedFile) (*graphql.TypesFile, error) {
	if file.Config.ModelsGoPath == "" {
		return nil, errors.Errorf("file: `%s`. Need to specify `models_go_path` option", file.Config.Name)
	}
	inputs, err := p.fileInputObjects(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file input objects")
	}
	inputsResolvers, err := p.fileInputMessagesResolvers(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file messages resolvers")
	}
	mapInputs, err := p.fileMapInputMessages(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map input objects")
	}
	mapOutputs, err := p.fileMapOutputMessages(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map output objects")
	}
	mapResolvers, err := p.fileInputMapResolvers(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file map resolvers")
	}
	outputMessages, err := p.fileOutputMessages(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file output messages")
	}
	services, err := p.fileServices(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare file services")
	}
	res := &graphql.TypesFile{
		PackageName:             file.OutputPkgName,
		Package:                 file.OutputPkg,
		InputObjects:            inputs,
		InputObjectResolvers:    inputsResolvers,
		OutputObjects:           outputMessages,
		MapInputObjects:         mapInputs,
		MapInputObjectResolvers: mapResolvers,
		MapOutputObjects:        mapOutputs,
		Services:                services,
	}

	return res, nil
}

func (p *Plugin) Prepare() error {
	parser := parser.Parser{}
	for _, cfg := range p.config.Files {
		file, err := os.Open(cfg.Path)
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		pf, err := parser.Parse(cfg.Path, file)
		if err != nil {
			file.Close()

			return errors.Wrap(err, "failed to parse swagger config")
		}
		err = file.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close file")
		}
		outPath, err := p.fileOutputPath(cfg)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve cfg '%s' output path", cfg.Path)
		}
		outPkgName, outPkg, err := p.fileOutputPackage(cfg)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve cfg '%s' output Go package", cfg.Path)
		}
		f := &parsedFile{
			File:          pf,
			Config:        cfg,
			OutputPath:    outPath,
			OutputPkg:     outPkg,
			OutputPkgName: outPkgName,
		}
		gqlFile, err := p.prepareTypesFile(f)
		if err != nil {
			return errors.Wrap(err, "failed to prepare types cfg")
		}
		p.graphql.AddTypesFile(outPath, gqlFile)
	}

	return nil
}

func (Plugin) Name() string {
	return PluginName
}

func (Plugin) Generate() error {
	return nil
}
