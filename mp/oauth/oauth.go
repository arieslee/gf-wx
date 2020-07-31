/*
@Time : 2020/7/30 5:58 下午
@Author : sunmoon
@File : oauth
@Software: GoLand
*/
package oauth

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"net/url"
)

const (
	authorizeURL          = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	accessTokenURL        = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&&secret=%s&code=%s&grant_type=authorization_code"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	accessTokenCacheKey   = "gf-wx-oauth-access-token"
)

type Oauth struct {
	config *config.MpConfig
}
type AccessTokenResult struct {
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	// UnionID 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	// 公众号文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
	UnionID string `json:"unionid"`
}

//UserInfo 用户授权获取到用户信息
type UserInfo struct {
	ErrCode    int64    `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// AuthCodeURL 生成网页授权地址.
//  appId:       公众号的唯一标识
//  redirectURI: 授权后重定向的回调链接地址
//  scope:       应用授权作用域 snsapi_userinfo or snsapi_base
//  state:       重定向后会带上 state 参数, 开发者可以填写 a-zA-Z0-9 的参数值, 最多128字节
func (o *Oauth) GetAuthCodeURL(redirectURI, scope, state string) string {
	return fmt.Sprintf(authorizeURL, url.QueryEscape(o.config.AppID), url.QueryEscape(redirectURI), url.QueryEscape(scope), url.QueryEscape(state))
}

// 获取accessToken
func (o *Oauth) GetAccessToken(code string) (*AccessTokenResult, error) {
	key := accessTokenCacheKey
	accessToken, _ := redis.String(g.Redis().DoVar("GET", key))
	result := &AccessTokenResult{}
	if len(accessToken) <= 0 {
		getTokenUrl := fmt.Sprintf(accessTokenURL, o.config.AppID, o.config.AppSecret, code)
		response := ghttp.GetBytes(getTokenUrl)
		err := gjson.DecodeTo(response, &result)
		if err != nil {
			glog.Line().Fatalf("GetAccessToken报文解析失败，error : %v", err)
			return nil, errors.New(fmt.Sprintf("GetAccessToken报文解析失败，error : %v", err))
		}
		if result.ErrCode != 0 {
			glog.Line().Fatalf("GetUserAccessToken error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
			return nil, errors.New(fmt.Sprintf("GetUserAccessToken error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
		}
		value := gconv.Map(result)
		expire := result.ExpiresIn - 100
		g.Redis().Do("SETEX", key, expire, gconv.String(value))
		return result, nil
	} else {
		err := gjson.DecodeTo(accessToken, &result)
		if err != nil {
			glog.Line().Fatalf("缓存内容解析失败，error : %v", err)
			return nil, errors.New(fmt.Sprintf("缓存内容解析失败，error : %v", err))
		}
		return result, nil
	}
}

// 刷新accessToken
func (o *Oauth) RefreshAccessToken(refreshToken string) (*AccessTokenResult, error) {
	refreshTokenURL := fmt.Sprintf(refreshAccessTokenURL, o.config.AppID, refreshToken)
	response := ghttp.GetBytes(refreshTokenURL)
	result := &AccessTokenResult{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("RefreshAccessToken报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("RefreshAccessToken报文解析失败，error : %v", err))
	}
	return result, nil
}

// 获取用户信息
func (o *Oauth) GetUserInfo(accessToken, openID string) (*UserInfo, error) {
	getUserInfoURL := fmt.Sprintf(userInfoURL, accessToken, openID)
	response := ghttp.GetBytes(getUserInfoURL)
	result := &UserInfo{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetUserInfo报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("GetUserInfo报文解析失败，error : %v", err))
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("accessToken error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, errors.New(fmt.Sprintf("accessToken error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
	}
	return result, nil
}
