package generator

type PluginsConfigs map[string]interface{}
type GenerateConfig struct {
	GenerateTraces bool   `yaml:"generate_tracer"`
	VendorPath     string `yaml:"vendor_path"`
	PluginsConfigs `yaml:",inline"`
}
