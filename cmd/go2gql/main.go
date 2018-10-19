package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/EGT-Ukraine/go2gql/generator"
)

var appFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Value: "generate.yml",
	},
}

func main() {
	app := cli.App{
		Flags: appFlags,
		Before: func(c *cli.Context) error {
			cfg, err := ioutil.ReadFile(c.String("config"))
			if err != nil {
				return errors.Wrap(err, "Failed to read config file")
			}
			gc := new(generator.GenerateConfig)
			if err = yaml.Unmarshal(cfg, gc); err != nil {
				return errors.Wrap(err, "Failed to unmarshal config file")
			}

			if err = gc.ParseImports(); err != nil {
				return errors.Wrap(err, "Failed to parse config file imports")
			}

			g := &generator.Generator{
				Config: gc,
			}

			for _, plugin := range Plugins(c) {
				if err := g.RegisterPlugin(plugin); err != nil {
					return errors.Wrap(err, "Failed to register plugin")
				}
			}
			if err = g.Init(); err != nil {
				return errors.Wrap(err, "failed to initialize generator")
			}
			if err = g.Prepare(); err != nil {
				return errors.Wrap(err, "failed to prepare generator")
			}
			c.App.Metadata["generator"] = g

			return nil
		},
		Commands: []cli.Command{
			{
				Name:    "info-keys",
				Aliases: []string{"ik"},
				Usage:   "Print all possible info keys",
				Action: func(c *cli.Context) {
					g := c.App.Metadata["generator"].(*generator.Generator)
					for plugin, keys := range g.GetPluginsInfosKeys() {
						if len(keys) > 0 {
							fmt.Println(plugin)
							for key, description := range keys {
								fmt.Println("\t- "+key, "\t\t"+description)
							}
						}
					}
				},
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Print info",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name: "infos",
					},
				},
				Action: func(c *cli.Context) {
					g := c.App.Metadata["generator"].(*generator.Generator)
					g.PrintInfos(c.StringSlice("infos"))
				},
			},
		},
		Action: func(c *cli.Context) {
			g := c.App.Metadata["generator"].(*generator.Generator)
			err := g.Generate()
			if err != nil {
				panic(errors.Wrap(err, "failed to generate"))
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
