// +build NOT (linux OR darwin)

package main

import (
	"github.com/urfave/cli"

	"github.com/EGT-Ukraine/go2gql/generator"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/graphql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/proto2gql"
	"github.com/EGT-Ukraine/go2gql/generator/plugins/swagger2gql"
)

func Plugins(c *cli.Context) []generator.Plugin {
	return []generator.Plugin{
		new(graphql.Plugin),
		new(swagger2gql.Plugin),
		new(proto2gql.Plugin),
	}
}
