package jumpServerSvr

import (
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl"
	"github.com/golang-jwt/jwt/v4"
)

type JumpServerService struct {
}

var jwtSecret = []byte("YOUR_SUPER_SECRET")

func (s *JumpServerService) GenerateJumpToken(userID, serverID, jumpType int) (string, error) {
	var addr = "http://127.0.0.1:8080/api/jumpServer/getServer"
	if len(global.GVA_CONFIG.JumpServer.OutAddr) > 0 && global.GVA_CONFIG.JumpServer.WebPort > 0 {
		addr = "http://" + global.GVA_CONFIG.JumpServer.OutAddr + ":" + strconv.Itoa(global.GVA_CONFIG.JumpServer.WebPort) + "/api/jumpServer/getServer"
	}
	claims := jumpServerMdl.JumpClaims{
		UserID:   userID,
		TargetID: serverID,
		JumpType: jumpType,
		Addr:     addr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)), // 1分钟有效
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *JumpServerService) ParseJumpToken(tokenStr string) (*jumpServerMdl.JumpClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jumpServerMdl.JumpClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jumpServerMdl.JumpClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
