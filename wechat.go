/*
@Time : 2020/7/30 6:06 下午
@Author : sunmoon
@File : wechat
@Software: GoLand
*/
package gf_wx

import (
    "github.com/arieslee/gf-wx/mp/config"
)

// Wechat struct
type Wechat struct {
}

// NewWechat init
func NewWechat() *Wechat {
    return &Wechat{}
}

func (w *Wechat) GetMp(cfg *config.MpConfig) {

}