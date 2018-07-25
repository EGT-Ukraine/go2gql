package graphql

type SchemaNodeConfig struct {
	Type           string             `mapstructure:"type"` // "OBJECT|SERVICE"
	Service        string             `mapstructure:"service"`
	ObjectName     string             `mapstructure:"object_name"`
	Field          string             `mapstructure:"field"`
	Fields         []SchemaNodeConfig `mapstructure:"fields"`
	ExcludeMethods []string           `mapstructure:"exclude_methods"`
	FilterMethods  []string           `mapstructure:"filter_methods"`
}
type SchemaConfig struct {
	Name          string            `mapstructure:"name"`
	OutputPath    string            `mapstructure:"output_path"`
	OutputPackage string            `mapstructure:"output_package"`
	Queries       *SchemaNodeConfig `mapstructure:"queries"`
	Mutations     *SchemaNodeConfig `mapstructure:"mutations"`
}
