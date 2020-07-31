/*
@Time : 2020/7/31 2:30 下午
@Author : sunmoon
@File : main
@Software: GoLand
*/
package main

import (
	gf_wx "github.com/arieslee/gf-wx"
	"github.com/arieslee/gf-wx/mp/config"
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
		state := r.GetString("state")
		glog.Line().Println("code", code)
		glog.Line().Println("state", state)
		oauth.GetAccessToken(code)
	})
	s.SetPort(8080)
	s.Run()
}
