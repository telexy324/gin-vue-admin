package config

import "os"

type Ansible struct {
	TmpPath             string `mapstructure:"tmp-path" json:"tmp-path" yaml:"tmp-path"`
	AccessKeyEncryption string `mapstructure:"access_key_encryption" json:"access_key_encryption" yaml:"access_key_encryption"`
	// task concurrency
	MaxParallelTasks int `mapstructure:"max_parallel_tasks" json:"max_parallel_tasks" yaml:"max_parallel_tasks"`
	// email alerting
	EmailSender   string `mapstructure:"email_sender" json:"email_sender" yaml:"email_sender"`
	EmailHost     string `mapstructure:"email_host" json:"email_host" yaml:"email_host"`
	EmailPort     string `mapstructure:"email_port" json:"email_port" yaml:"email_port"`
	EmailUsername string `mapstructure:"email_username" json:"email_username" yaml:"email_username"`
	EmailPassword string `mapstructure:"email_password" json:"email_password" yaml:"email_password"`
	// telegram alerting
	TelegramChat  string `mapstructure:"telegram_chat" json:"telegram_chat" yaml:"telegram_chat"`
	TelegramToken string `mapstructure:"telegram_token" json:"telegram_token" yaml:"telegram_token"`
	// feature switches
	EmailAlert    bool `mapstructure:"email_alert" json:"email_alert" yaml:"email_alert"`
	EmailSecure   bool `mapstructure:"email_secure" json:"email_secure" yaml:"email_secure"`
	TelegramAlert bool `mapstructure:"telegram_alert" json:"telegram_alert" yaml:"telegram_alert"`
	// web host
	WebHost string `mapstructure:"web_host" json:"web_host" yaml:"web_host"`
}

func (conf *Ansible) GetAccessKeyEncryption() string {
	ret := os.Getenv("SEMAPHORE_ACCESS_KEY_ENCRYPTION")

	if ret == "" {
		ret = conf.AccessKeyEncryption
	}

	return ret
}
