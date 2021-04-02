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
	mp_server "github.com/arieslee/gf-wx/mp/server"
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
func (m *Mp) NewServer()*mp_server.MPServer {
	return mp_server.NewMPServer(m.config)
}