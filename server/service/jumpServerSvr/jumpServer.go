package jumpServerSvr

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/jumpServer"
)

type JumpServerService struct {
}

//var jwtSecret = []byte("YOUR_SUPER_SECRET")
//
//func (s *JumpServerService) GenerateJumpToken(userID, serverID, jumpType int) (string, error) {
//	var addr = "http://127.0.0.1:8080/api/jumpServer/getServer"
//	if len(global.GVA_CONFIG.JumpServer.OutAddr) > 0 && global.GVA_CONFIG.JumpServer.WebPort > 0 {
//		addr = "http://" + global.GVA_CONFIG.JumpServer.OutAddr + ":" + strconv.Itoa(global.GVA_CONFIG.JumpServer.WebPort) + "/api/jumpServer/getServer"
//	}
//	claims := jumpServerMdl.JumpClaims{
//		UserID:   userID,
//		TargetID: serverID,
//		JumpType: jumpType,
//		Addr:     addr,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)), // 1分钟有效
//			IssuedAt:  jwt.NewNumericDate(time.Now()),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(jwtSecret)
//}
//
//func (s *JumpServerService) ParseJumpToken(tokenStr string) (*jumpServerMdl.JumpClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenStr, &jumpServerMdl.JumpClaims{}, func(t *jwt.Token) (interface{}, error) {
//		return jwtSecret, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := token.Claims.(*jumpServerMdl.JumpClaims)
//	if !ok || !token.Valid {
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	return claims, nil
//}

func genSecret() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b) // 无特殊字符
}

var hmacKey = []byte("bastion-super-secret-key")

func (s *JumpServerService) GenerateSession(server application.ApplicationServer, client string) (string, *jumpServerMdl.SessionPayload, error) {
	bastionHost := "127.0.0.1"
	if len(global.GVA_CONFIG.JumpServer.OutAddr) > 0 {
		bastionHost = global.GVA_CONFIG.JumpServer.OutAddr
	}
	payload := jumpServerMdl.SessionPayload{
		BastionHost: bastionHost,
		BastionPort: global.GVA_CONFIG.JumpServer.Port,
		Client:      client,
		Secret:      genSecret(),
		IssuedAt:    time.Now().Unix(),
		ExpireAt:    time.Now().Add(24 * time.Hour).Unix(),
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(raw)

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(payloadB64))
	sig := mac.Sum(nil)

	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	token := payloadB64 + "." + sigB64
	sessRec := jumpServerMdl.SessionRecord{
		Secret:     token,
		UserID:     0,
		TargetHost: server.ManageIp,
		TargetPort: server.SshPort,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		Used:       false,
	}
	if err = jumpServer.SessStore.Save(payload.Secret, &sessRec); err != nil {
		return "", nil, err
	}
	return token, &payload, nil
}
