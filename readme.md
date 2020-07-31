```
package main

import (
	"fmt"
	"github.com/arieslee/gf-wx"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
)

func main() {
	s := g.Server()
	cfg := &config.MpConfig{
		AppID:     "123",
		AppSecret: "123",
	}
	wechat := gf_wx.NewWechat()
	wechat.InitRedis() //如果需要的话
	mp := wechat.GetMp(cfg)
	oauth := mp.GetOauth()
	s.BindHandler("GET:/index", func(r *ghttp.Request) {
		state := grand.S(32)
		g.Redis().Do("SET", "state", state)
		redirectURL := oauth.GetAuthCodeURL("http://wx-test.hchmc.cn/callback", "snsapi_userinfo", state)
		r.Response.RedirectTo(redirectURL)
	})
	// 获取微信用户信息
	s.BindHandler("GET:/callback", func(r *ghttp.Request) {
		code := r.GetString("code")
		state := r.GetString("state")
		localState, _ := g.Redis().Do("GET", "state")
		if gconv.String(localState) != state {
			fmt.Printf("非法请求")
			return
		}
		accessToken, err := oauth.GetAccessToken(code)
		if err != nil {
			fmt.Println(err.Error())
		}
		if accessToken != nil {
			openId := accessToken.OpenID
			token := accessToken.AccessToken
			info, _ := oauth.GetUserInfo(token, openId)
			glog.Line().Println(info)
		} else {
			fmt.Printf("accessToken无效")
			return
		}
	})
	s.SetPort(8080)
	s.Run()
}
```