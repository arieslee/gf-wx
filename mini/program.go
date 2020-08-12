/*
@Time : 2020/8/12 2:57 下午
@Author : sunmoon
@File : program
@Software: GoLand
*/
package mini

import (
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/arieslee/gf-wx/mini/program/login"
	"github.com/arieslee/gf-wx/mini/program/payment"
	"github.com/arieslee/gf-wx/mini/program/qrcode"
)

type Program struct {
	config *config.MiniConfig
}

func NewProgram(cfg *config.MiniConfig) *Program {
	return &Program{
		config: cfg,
	}
}

func (p *Program) GetLogin() *login.Login {
	return login.NewLogin(p.config)
}

func (p *Program) GetQrcode() *qrcode.QRCode {
	return qrcode.NewQRCode(p.config)
}
func (p *Program) GetPayment() *payment.Payment {
	return payment.NewPayment(p.config)
}
