/*
@Time : 2020/7/30 6:06 下午
@Author : sunmoon
@File : wechat
@Software: GoLand
*/
package gf_wx

import (
	"github.com/arieslee/gf-wx/mp"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/gogf/gf/database/gredis"
)

// Wechat struct
type Wechat struct {
}

// NewWechat init
func NewWechat() *Wechat {
	return &Wechat{}
}

func (w *Wechat) InitRedis() {
	// 初始化gredis
	redisCfg := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}
	gredis.SetConfig(redisCfg)
}

func (w *Wechat) GetMp(cfg *config.MpConfig) *mp.Mp {
	return mp.NewMP(cfg)
}
