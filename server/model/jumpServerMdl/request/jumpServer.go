package request

type GetToken struct {
	ID   float64 `json:"id" form:"id"` // 主键ID
	Type int     `json:"type" form:"type"`
}

type Token struct {
	Token string `json:"token"` // token
}
