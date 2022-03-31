package config

import "os"

type Ansible struct {
	TmpPath string `mapstructure:"tmp-path" json:"tmp-path" yaml:"tmp-path"`
	AccessKeyEncryption string `mapstructure:"access_key_encryption" json:"access_key_encryption" yaml:"access_key_encryption"`
}

func (conf *Ansible) GetAccessKeyEncryption() string {
	ret := os.Getenv("SEMAPHORE_ACCESS_KEY_ENCRYPTION")

	if ret == "" {
		ret = conf.AccessKeyEncryption
	}

	return ret
}