/*
@Time : 2020/7/31 2:30 下午
@Author : sunmoon
@File : main
@Software: GoLand
*/
package main

import (
	"github.com/arieslee/gf-wx"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/grand"
)

func main() {
	s := g.Server()
	cfg := &config.MpConfig{
		AppID:     "123",
		AppSecret: "123",
	}
	var (
		redisCfg = gredis.Config{
			Host: "127.0.0.1",
			Port: 6379,
			Db:   1,
		}
	)
	gredis.SetConfig(redisCfg)
	wechat := gf_wx.NewWechat()
	mp := wechat.GetMp(cfg)
	oauth := mp.GetOauth()
	s.BindHandler("GET:/index", func(r *ghttp.Request) {
		state := grand.S(32)
		redirectURL := oauth.GetAuthCodeURL("http://wx-test.hchmc.cn/callback", "snsapi_userinfo", state)
		r.Response.RedirectTo(redirectURL)
	})
	s.BindHandler("GET:/callback", func(r *ghttp.Request) {
		code := r.GetString("code")
		accessToken, err := oauth.GetAccessToken(code)
		if err != nil {
			glog.Line().Println(err.Error())
		}
		openId := accessToken.OpenID
		token := accessToken.AccessToken
		info, _ := oauth.GetUserInfo(token, openId)
		glog.Line().Println(info)
	})
	s.SetPort(8080)
	s.Run()
}
