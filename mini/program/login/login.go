/*
@Time : 2020/8/12 12:59 下午
@Author : sunmoon
@File : oauth
@Software: GoLand
*/
package login

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	codeToSessionUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type MiniProgramLogin struct {
	config *config.MiniConfig
}

type CodeToSessionResult struct {
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	// UnionID 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	// 公众号文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
	UnionID string `json:"unionid"`
}

func NewLogin(cfg *config.MiniConfig) *MiniProgramLogin {
	return &MiniProgramLogin{
		config: cfg,
	}
}

// CodeToSession 登录凭证校验
// @see https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (o *MiniProgramLogin) CodeToSession(jsCode string) (*CodeToSessionResult, error) {
	URL := fmt.Sprintf(codeToSessionUrl, o.config.AppID, o.config.AppSecret, jsCode)
	response := ghttp.GetBytes(URL)
	result := &CodeToSessionResult{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		g.Log().Line().Fatalf("CodeToSession报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("CodeToSession报文解析失败，error : %v", err))
	}
	if result.ErrCode != 0 {
		g.Log().Line().Fatalf("CodeToSession error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, errors.New(fmt.Sprintf("CodeToSession error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
	}
	return result, nil
}
