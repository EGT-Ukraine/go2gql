package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql"

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
			cfgFile, err := os.Open(c.String("config"))
			if err != nil {
				panic(err)
			}
			cfg, err := ioutil.ReadAll(cfgFile)
			if err != nil {
				panic(err)
			}
			gc := new(generator.GenerateConfig)
			err = yaml.Unmarshal(cfg, gc)
			if err != nil {
				panic(err)
			}
			g := &generator.Generator{
				Config: gc,
			}

			for _, plugin := range []generator.Plugin{
				new(graphql.Plugin),
				new(proto2gql.Plugin),
				new(swagger2gql.Plugin),
			} {
				err := g.RegisterPlugin(plugin)
				if err != nil {
					panic(err.Error())
				}
			}
			err = g.Init()
			if err != nil {
				panic(errors.Wrap(err, "failed to initialize generator"))
			}
			err = g.Prepare()
			if err != nil {
				panic(errors.Wrap(err, "failed to prepare generator"))
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
	app.Run(os.Args)
}
