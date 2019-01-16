package dataloader

type DataLoadersConfig struct {
	OutputPath string `mapstructure:"output_path"`
}
type DataLoaderProviderConfig struct {
	Name           string `mapstructure:"name"`
	WaitDurationMs int    `mapstructure:"wait_duration_ms"`
}
type DataLoaderFieldConfig struct {
	FieldName     string `mapstructure:"field_name"`
	KeyFieldName  string `mapstructure:"key_field_name"`
	KeyFieldSlice bool   `mapstructure:"key_field_slice"`
	DataLoader    string `mapstructure:"data_loader_name"`
}
