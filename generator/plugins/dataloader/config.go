package dataloader

type DataLoadersConfig struct {
	OutputPath string `mapstructure:"output_path"`
}

type FieldConfig struct {
	FieldName    string `mapstructure:"field_name"`
	KeyFieldName string `mapstructure:"key_field_name"`
	DataLoader   string `mapstructure:"data_loader_name"`
}
