package main

import (
	"github.com/arieslee/gf-wx"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	serverName := g.Cfg().GetString("server.name")
	s := g.Server(serverName)
	appId := g.Cfg().GetString("wechat.appId")
	appSecret := g.Cfg().GetString("wechat.appSecret")
	token := g.Cfg().GetString("wechat.token")
	cfg := &config.MpConfig{
		AppID:    appId,
		AppSecret: appSecret,
		Token: token,
	}
	wechat := gf_wx.NewWechat()
	mp := wechat.GetMp(cfg)
	mpServer := mp.NewServer()
	s.BindHandler("ALL:/server", func(r *ghttp.Request) {
		mpServer.Build(r)
	})
	port := g.Cfg().GetInt("server.port")
	s.SetPort(port)
	s.Run()
}