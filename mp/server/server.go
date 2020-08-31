/*
@Time : 2020/8/14 11:37 上午
@Author : sunmoon
@File : server
@Software: GoLand
*/
package mp_server

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/helper"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/arieslee/gf-wx/mp/token"
	"github.com/gogf/gf/net/ghttp"
)

type MPServer struct {
	config *config.MpConfig
}

func NewMPServer(cfg *config.MpConfig) *MPServer {
	return &MPServer{
		config: cfg,
	}
}
func (s *MPServer) Build(r *ghttp.Request) error {
	encryptType := r.GetString("encrypt_type")
	switch encryptType {
	case "aes":
		haveSignature := r.GetString("signature")
		if len(haveSignature) == 0 {
			return errors.New("not found signature query parameter")
		}
		haveMsgSignature := r.GetString("msg_signature")
		if len(haveMsgSignature) == 0 {
			return errors.New("not found msg_signature query parameter")
		}
		timestampString := r.GetString("timestamp")
		if len(timestampString) == 0 {
			return errors.New("not found timestamp query parameter")
		}
		nonce := r.GetString("nonce")
		if len(nonce) == 0 {
			return errors.New("not found nonce query parameter")
		}
		tokenStr, err := s.GetToken()
		if err != nil {
			return err
		}
		wantSignature := helper.Sign(tokenStr, timestampString, nonce)
		if haveSignature != wantSignature {
			err = fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
		}
	}
	//@todo
	return nil
}
func (s *MPServer) GetToken() (string, error) {
	tokenServ := token.NewToken(s.config)
	accessToken, err := tokenServ.GetToken()
	if err != nil {
		return "", err
	}
	return accessToken.AccessToken, nil
}
