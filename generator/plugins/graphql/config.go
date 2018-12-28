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
type DataLoadersConfig struct {
	OutputPath string             `mapstructure:"output_path"`
	Loaders    []DataLoaderConfig `mapstructure:"loaders"`
}
type DataLoaderConfig struct {
	Name           string `mapstructure:"name"`
	ServiceName    string `mapstructure:"service_name"`
	MethodName     string `mapstructure:"method_name"`
	WaitDurationMs int    `mapstructure:"wait_duration_ms"`
}
