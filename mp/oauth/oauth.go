/*
@Time : 2020/7/30 5:58 下午
@Author : sunmoon
@File : oauth
@Software: GoLand
*/
package oauth

import (
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "github.com/gogf/gf/os/glog"
    "github.com/gomodule/redigo/redis"
    "net/url"
)
const (
    authorizeURL="https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
    accessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&&secret=%s&code=%s&grant_type=authorization_code"
    AccessTokenCacheKey = "gf-wx-oauth-access-token"
)
type Oauth struct {

}
// AuthCodeURL 生成网页授权地址.
//  appId:       公众号的唯一标识
//  redirectURI: 授权后重定向的回调链接地址
//  scope:       应用授权作用域 snsapi_userinfo or snsapi_base
//  state:       重定向后会带上 state 参数, 开发者可以填写 a-zA-Z0-9 的参数值, 最多128字节
func GetAuthCodeURL(appId, redirectURI, scope, state string) string {
    return fmt.Sprintf(authorizeURL, url.QueryEscape(appId), url.QueryEscape(redirectURI), url.QueryEscape(scope), url.QueryEscape(state))
}

func GetAccessToken(code string)  {
    key := AccessTokenCacheKey
    accessToken, _ := redis.String(g.Redis().DoVar("GET", key))
    if len(accessToken) <= 0{
        getTokenUrl := fmt.Sprintf(accessTokenURL, "", "", code)
        response := ghttp.GetBytes(getTokenUrl)
        glog.Line().Println(response)
    }
}