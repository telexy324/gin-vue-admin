package config

type Ansible struct {
	TmpPath string `mapstructure:"tmp-path" json:"tmp-path" yaml:"tmp-path"`
}
