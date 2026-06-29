package config

type JumpServer struct {
	Port           int    `mapstructure:"port" json:"port" yaml:"port"` // 端口值
	WebPort        int    `mapstructure:"web-port" json:"webPort" yaml:"webPort"`
	WebSocketPath  string `mapstructure:"websocket-path" json:"websocketPath" yaml:"websocketPath"`
	WebSocketURL   string `mapstructure:"websocket-url" json:"websocketUrl" yaml:"websocketUrl"`
	IdleTimeout    int    `mapstructure:"idle-time-out" json:"idleTimeout" yaml:"idleTimeout"`
	MaxSessionTime int    `mapstructure:"max-session-time" json:"maxSessionTime" yaml:"maxSessionTime"`
	OutAddr        string `mapstructure:"out-addr" json:"outAddr" yaml:"outAddr"`
}
