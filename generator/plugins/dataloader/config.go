package dataloader

import "time"

type DataLoadersConfig struct {
	OutputPath string `mapstructure:"output_path"`
}

type ProviderConfig struct {
	Name         string        `mapstructure:"name"`
	WaitDuration time.Duration `mapstructure:"wait_duration"`
	Slice        bool
}

type FieldConfig struct {
	FieldName    string `mapstructure:"field_name"`
	KeyFieldName string `mapstructure:"key_field_name"`
	DataLoader   string `mapstructure:"data_loader_name"`
}
