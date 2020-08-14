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
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

const (
	codeToSessionUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=JSCODE&grant_type=authorization_code"
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

func (o *MiniProgramLogin) CodeToSession() (*CodeToSessionResult, error) {
	URL := fmt.Sprintf(codeToSessionUrl, o.config.AppID, o.config.AppSecret)
	response := ghttp.GetBytes(URL)
	result := &CodeToSessionResult{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("CodeToSession报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("CodeToSession报文解析失败，error : %v", err))
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("CodeToSession error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, errors.New(fmt.Sprintf("CodeToSession error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
	}
	return result, nil
}
