package generator

import (
	"github.com/pkg/errors"
)

type Infos []string

func (i Infos) Contains(v string) bool {
	for _, val := range i {
		if val == v {
			return true
		}
	}
	return false
}

type Plugin interface {
	Init(*GenerateConfig, []Plugin) error
	Prepare() error
	Name() string
	PrintInfo(info Infos)
	Infos() map[string]string
	Generate() error
}
type PluginFabric = func() Plugin

type Generator struct {
	Config  *GenerateConfig
	Plugins []Plugin
}

type PluginContext struct {
	Plugins []Plugin
}

func (g *Generator) RegisterPlugin(p Plugin) error {
	for _, plugin := range g.Plugins {
		if plugin.Name() == p.Name() {
			return errors.Errorf("plugin with name '%s' already exists", p.Name())
		}
	}
	g.Plugins = append(g.Plugins, p)
	return nil
}
func (g *Generator) Init() error {
	for _, plugin := range g.Plugins {
		err := plugin.Init(g.Config, g.Plugins)
		if err != nil {
			return errors.Wrapf(err, "failed to initialize plugin %s", plugin.Name())
		}
	}
	return nil
}

func (g *Generator) Prepare() error {
	for _, plugin := range g.Plugins {
		err := plugin.Prepare()
		if err != nil {
			return errors.Wrapf(err, "failed to prepare plugin %s", plugin.Name())
		}
	}
	return nil
}

func (g *Generator) Generate() error {
	for _, plugin := range g.Plugins {
		err := plugin.Generate()
		if err != nil {
			return errors.Wrapf(err, "plugin %s generation errors", plugin.Name())
		}
	}
	return nil
}
func (g *Generator) PrintInfos(i []string) {
	for _, p := range g.Plugins {
		p.PrintInfo(Infos(i))
	}
}
func (g *Generator) GetPluginsInfosKeys() map[string]map[string]string {
	res := make(map[string]map[string]string)
	for _, p := range g.Plugins {
		res[p.Name()] = p.Infos()
	}
	return res
}
