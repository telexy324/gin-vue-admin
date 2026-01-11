package jumpServerSvr

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JumpServerService struct {
}

var jwtSecret = []byte("YOUR_SUPER_SECRET")

func (s *JumpServerService) GenerateJumpToken(userID, serverID int) (string, error) {
	claims := jumpServerMdl.JumpClaims{
		UserID:   userID,
		TargetID: serverID,
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
