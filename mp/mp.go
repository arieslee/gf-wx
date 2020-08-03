/*
@Time : 2020/7/30 6:19 下午
@Author : sunmoon
@File : mp
@Software: GoLand
*/
package mp

import (
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/arieslee/gf-wx/mp/js"
	"github.com/arieslee/gf-wx/mp/oauth"
	"github.com/arieslee/gf-wx/mp/server"
	"github.com/gogf/gf/net/ghttp"
)

type Mp struct {
	config *config.MpConfig
}

func NewMP(cfg *config.MpConfig) *Mp {
	return &Mp{
		config: cfg,
	}
}

func (m *Mp) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(m.config)
}
func (m *Mp) GetJS() *js.Js {
	return js.NewJS(m.config)
}
func (m *Mp) GetServer(req *ghttp.Request) *server.Server {
	return server.NewServer(req.Response, req, m.config.AppID, m.config.Token, m.config.EncodingAESKey)
}
