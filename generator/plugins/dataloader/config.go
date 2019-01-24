package dataloader

type DataLoadersConfig struct {
	OutputPath string `mapstructure:"output_path"`
}
type ProviderConfig struct {
	Name           string `mapstructure:"name"`
	WaitDurationMs int    `mapstructure:"wait_duration_ms"`
}
type FieldConfig struct {
	FieldName    string `mapstructure:"field_name"`
	KeyFieldName string `mapstructure:"key_field_name"`
	DataLoader   string `mapstructure:"data_loader_name"`
}
