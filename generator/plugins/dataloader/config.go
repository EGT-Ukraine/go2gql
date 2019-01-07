package dataloader

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
