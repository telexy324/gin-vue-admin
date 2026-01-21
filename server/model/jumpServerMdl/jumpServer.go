package jumpServerMdl

import "github.com/golang-jwt/jwt/v4"

type JumpClaims struct {
	UserID   int    `json:"user_id"`
	TargetID int    `json:"target_id"`
	JumpType int    `json:"jump_type"`
	Addr     string `json:"addr"`
	jwt.RegisteredClaims
}

type ConnInfo struct {
	JumpHost string `json:"jump_host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Protocol string `json:"protocol"` // ssh / sftp
	Client   string `json:"client"`   // securecrt / filezilla
	Password string `json:"password"` // 临时凭证（一次性）
}
