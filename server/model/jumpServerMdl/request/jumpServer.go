package request

type GetToken struct {
	ID     float64 `json:"id" form:"id"` // 主键ID
	Client string  `json:"client" form:"client"`
}

type Token struct {
	Token string `json:"token"` // token
}
