package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ImportedPluginsConfigs struct {
	Path           string
	PluginsConfigs PluginsConfigs
}

type PluginsConfigs map[string]interface{}
type GenerateConfig struct {
	GenerateTraces        bool     `yaml:"generate_tracer"`
	VendorPath            string   `yaml:"vendor_path"`
	Imports               []string `yaml:"imports"`
	PluginsConfigsImports []ImportedPluginsConfigs
	PluginsConfigs        `yaml:",inline"`
}

func (gc *GenerateConfig) ParseImports() error {
	for _, importPath := range gc.Imports {
		normalizedPath := os.ExpandEnv(importPath)
		normalizedPath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return errors.Wrapf(err, "failed to make normalized path '%s' absolute", normalizedPath)
		}

		cfg, err := ioutil.ReadFile(normalizedPath)
		if err != nil {
			return errors.Wrapf(err, "Failed to read import '%s' file", normalizedPath)
		}

		pluginsConfig := PluginsConfigs{}

		importedPluginsConfig := ImportedPluginsConfigs{
			Path:           normalizedPath,
			PluginsConfigs: pluginsConfig,
		}

		err = yaml.Unmarshal(cfg, pluginsConfig)
		if err != nil {
			return errors.Wrapf(err, "Failed to unmarshal import '%s' file", normalizedPath)
		}

		gc.PluginsConfigsImports = append(gc.PluginsConfigsImports, importedPluginsConfig)
	}

	return nil
}
