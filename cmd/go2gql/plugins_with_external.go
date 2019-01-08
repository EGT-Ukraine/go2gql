// +build linux darwin

package main

import (
	"path/filepath"
	"plugin"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/EGT-Ukraine/go2gql/generator"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/dataloader"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql"
)

func Plugins(c *cli.Context) []generator.Plugin {
	res := []generator.Plugin{
		new(dataloader.Plugin),
		new(graphql.Plugin),
		new(swagger2gql.Plugin),
		new(proto2gql.Plugin),
	}
	pluginsDir := c.String("plugins")
	if len(pluginsDir) > 0 {
		plugins, err := filepath.Glob(filepath.Join(pluginsDir, "*.so"))
		if err != nil {
			panic(errors.Wrap(err, "failed to scan plugins directory"))
		}
		for _, pluginPath := range plugins {
			p, err := plugin.Open(pluginPath)
			if err != nil {
				panic(errors.Wrapf(err, "failed to open plugin %s", pluginPath))
			}
			s, err := p.Lookup("Plugin")
			if err != nil {
				panic(errors.Wrapf(err, "no `Plugin` symbol in plugin %s", pluginPath))
			}
			pf, ok := s.(generator.PluginFabric)
			if !ok {
				panic(errors.Errorf("symbol `Plugin` does not implements `func() Plugin` interface in plugin %s", pluginPath))
			}
			res = append(res, pf())
		}
	}

	return res
}

func init() {
	appFlags = append(appFlags, cli.StringFlag{
		Name:  "plugins, p",
		Usage: "Plugins directory",
	})
}
