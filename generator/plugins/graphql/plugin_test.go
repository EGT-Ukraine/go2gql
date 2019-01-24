package graphql

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	yaml "gopkg.in/yaml.v2"

	"github.com/EGT-Ukraine/go2gql/generator"
)

func TestSchemaMergeRecursively(t *testing.T) {
	Convey("Given config with import", t, func() {
		var mainConfig = `
graphql_schemas:
  - name: "API"
    output_path: "./services_api/schema/api.go"
    output_package: "schema"
    queries:
      type: "OBJECT"
      fields:
        - field: "common"
          object_name: "Common"
          type: "OBJECT"
          fields:
            - field: "commonNested"
              object_name: "CommonNested"
              type: "OBJECT"
              fields:
                - field: "service1Field"
                  object_name: "Service1Object"
                  service: "Service1"
                  type: "SERVICE"
`

		var importConfig = `
graphql_schemas:
  - name: "API"
    queries:
      type: "OBJECT"
      fields:
        - field: "common"
          object_name: "Common"
          type: "OBJECT"
          fields:
            - field: "commonNested"
              object_name: "CommonNested"
              type: "OBJECT"
              fields:
                - field: "service2Field"
                  object_name: "Service2Object"
                  service: "Service2"
                  type: "SERVICE"
`
		gc, err := parseConfigs(mainConfig, importConfig)

		So(err, ShouldBeNil)

		Convey("When the graphql plugin is initialized", func() {
			graphqlPlugin := new(Plugin)

			if err := graphqlPlugin.Init(gc, []generator.Plugin{}); err != nil {
				t.Fatalf(err.Error())
			}

			Convey("Graphql schema should be merged recursively", func() {
				So(graphqlPlugin.schemaConfigs[0], ShouldResemble, SchemaConfig{
					Name:          "API",
					OutputPath:    "./services_api/schema/api.go",
					OutputPackage: "schema",
					Queries: &SchemaNodeConfig{
						Type: "OBJECT",
						Fields: []SchemaNodeConfig{
							{
								Type:       "OBJECT",
								ObjectName: "Common",
								Field:      "common",
								Fields: []SchemaNodeConfig{
									{
										Type:       "OBJECT",
										ObjectName: "CommonNested",
										Field:      "commonNested",
										Fields: []SchemaNodeConfig{
											{
												Type:       "SERVICE",
												Service:    "Service1",
												ObjectName: "Service1Object",
												Field:      "service1Field",
											},
											{
												Type:       "SERVICE",
												Service:    "Service2",
												ObjectName: "Service2Object",
												Field:      "service2Field",
											},
										},
									},
								},
							},
						},
					},
				})
			})
		})
	})
}

func parseConfigs(mainConfig string, importConfig string) (*generator.GenerateConfig, error) {
	gc := new(generator.GenerateConfig)

	pluginsConfig := generator.PluginsConfigs{}

	importedPluginsConfig := generator.ImportedPluginsConfigs{
		Path:           "/home/go/import.yml",
		PluginsConfigs: pluginsConfig,
	}

	if err := yaml.Unmarshal([]byte(importConfig), pluginsConfig); err != nil {
		return nil, errors.Wrap(err, "Failed to parse import config")
	}

	gc.PluginsConfigsImports = append(gc.PluginsConfigsImports, importedPluginsConfig)

	if err := yaml.Unmarshal([]byte(mainConfig), gc); err != nil {
		return nil, errors.Wrap(err, "Failed to parse main config")
	}

	return gc, nil
}
