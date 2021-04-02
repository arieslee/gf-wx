/*
@Time : 2020/8/14 11:37 上午
@Author : sunmoon
@File : server
@Software: GoLand
*/
package mp_server

import (
	"crypto/sha1"
	"fmt"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/arieslee/gf-wx/mp/token"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"sort"
)

const (
	LangZhCN = "zh_CN" // 简体中文
	LangZhTW = "zh_TW" // 繁体中文
	LangEN   = "en"    // 英文
)

const (
	SexUnknown = 0 // 未知
	SexMale    = 1 // 男性
	SexFemale  = 2 // 女性
)

type MPServer struct {
	config *config.MpConfig
}

func NewMPServer(cfg *config.MpConfig) *MPServer {
	return &MPServer{
		config: cfg,
	}
}
func (s *MPServer) Signature(signature, timestamp, nonce string) (string,bool) {
	strs := sort.StringSlice{s.config.Token, timestamp, nonce}
	sort.Strings(strs)
	str := ""

	for _, s := range strs {
		str += s
	}

	h := sha1.New()
	h.Write([]byte(str))

	signatureNow := fmt.Sprintf("%x", h.Sum(nil))
	if signature == signatureNow {
		return signatureNow, true
	}
	return signatureNow, false
}
func (s *MPServer) Build(r *ghttp.Request) {
	echoStr := r.GetString("echostr")
	if len(echoStr) > 0{
		timestampString := r.GetString("timestamp")
		if len(timestampString) == 0 {
			g.Log().Line(true).Println("not found timestamp query parameter")
			return
		}
		nonce := r.GetString("nonce")
		if len(nonce) == 0 {
			g.Log().Line(true).Println("not found nonce query parameter")
			return
		}
		signature := r.GetString("signature")
		wantSignature,yes := s.Signature(signature, timestampString, nonce)
		if !yes {
			g.Log().Line(true).Println(fmt.Sprintf("check signature failed, have: %s, want: %s", signature, wantSignature))
			return
		}

	}else{
		s.Handler(r)
	}
	return
}
func (s *MPServer) Handler(r *ghttp.Request) {
	g.Log().Line(true).Println(r.GetMap())
}
func (s *MPServer) GetToken() (string, error) {
	tokenServ := token.NewToken(s.config)
	accessToken, err := tokenServ.GetToken()
	if err != nil {
		return "", err
	}
	return accessToken.AccessToken, nil
}
