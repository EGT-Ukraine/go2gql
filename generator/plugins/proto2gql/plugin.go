package proto2gql

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/EGT-Ukraine/go2gql/generator"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql/lib/pluginconfig"
)

const (
	PluginName            = "proto2gql"
	PluginConfigKey       = "proto2gql"
	PluginImportConfigKey = "proto2gql_files"
)

type Plugin struct {
	graphql          *graphql.Plugin
	dataLoaderPlugin *dataloader.Plugin
	config           *Config
	generateConfig   *generator.GenerateConfig
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
		return errors.New("'graphql' plugin is not installed")
	}
	if p.dataLoaderPlugin == nil {
		return errors.New("'dataloader' plugin is not installed")
	}
	cfg := new(Config)
	err := pluginconfig.Decode(config.PluginsConfigs[PluginConfigKey], cfg)
	if err != nil {
		return errors.Wrap(err, "failed to decode config")
	}
	p.generateConfig = config
	p.config = cfg

	if err = p.parseImports(); err != nil {
		return errors.Wrap(err, "failed to decode imports")
	}

	err = p.normalizeGenerateConfigPaths()
	if err != nil {
		return errors.Wrap(err, "failed to normalize config paths")
	}

	return nil
}

func (p *Plugin) parseImports() error {
	for _, pluginsConfigsImports := range p.generateConfig.PluginsConfigsImports {
		configs := new([]*ProtoFileConfig)
		if err := mapstructure.Decode(pluginsConfigsImports.PluginsConfigs[PluginImportConfigKey], configs); err != nil {
			return errors.Wrap(err, "failed to decode config")
		}

		for _, config := range *configs {
			var importFileDir = filepath.Dir(pluginsConfigsImports.Path)

			var protoPath = filepath.Join(importFileDir, config.ProtoPath)

			config.ProtoPath = protoPath
			config.Paths = append(config.Paths, importFileDir)
			p.config.Files = append(p.config.Files, config)
		}
	}

	return nil
}

func (p *Plugin) normalizeGenerateConfigPaths() error {
	for i, path := range p.config.Paths {
		normalizedPath := os.ExpandEnv(path)
		normalizedPath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return errors.Wrapf(err, "failed to make normalized path '%s' absolute", normalizedPath)
		}
		p.config.Paths[i] = normalizedPath
	}
	for i, file := range p.config.Files {
		normalizedPath := os.ExpandEnv(file.ProtoPath)
		normalizedPath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return errors.Wrapf(err, "failed to make normalized path '%s' absolute", normalizedPath)
		}
		p.config.Files[i].ProtoPath = normalizedPath

	}

	return nil
}

func (p *Plugin) prepareFileConfig(fileCfg *ProtoFileConfig) {
	fileCfg.Paths = append(fileCfg.Paths, p.config.Paths...)

	fileCfg.ImportsAliases = append(fileCfg.ImportsAliases, p.config.ImportsAliases...)
}

func (p *Plugin) PrintInfo(info generator.Infos) {
}

func (p *Plugin) Infos() map[string]string {
	return nil
}

func (p *Plugin) Prepare() error {
	pr := new(Proto2GraphQL)
	pr.VendorPath = p.generateConfig.VendorPath
	pr.DataLoaderPlugin = p.dataLoaderPlugin
	pr.GenerateTracers = p.generateConfig.GenerateTraces
	pr.OutputPath = p.config.GetOutputPath()
	for _, file := range p.config.Files {
		p.prepareFileConfig(file)

		if err := pr.AddSourceByConfig(file); err != nil {
			return errors.Wrap(err, "failed to parse file "+file.ProtoPath)
		}
	}
	for _, file := range pr.parser.ParsedFiles() {
		pf, err := pr.parsedFile(file)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve parsed file of '%s'", file.FilePath)
		}

		commonFile, err := pr.prepareFile(pf)
		if err != nil {
			return errors.Wrap(err, "failed to prepare file for generation")
		}
		p.graphql.AddTypesFile(pf.OutputPath, commonFile)
	}

	return nil
}

func (Plugin) Name() string {
	return PluginName
}

func (Plugin) Generate() error {
	return nil
}
