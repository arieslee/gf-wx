/*
@Time : 2020/8/12 1:30 下午
@Author : sunmoon
@File : qrcode
@Software: GoLand
*/
package qrcode

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/arieslee/gf-wx/mini/program/token"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type QRCode struct {
	config *config.MiniConfig
}

// Color QRCode color
type Color struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

const (
	createWXAQRCodeURL   = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s"
	getWXACodeURL        = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	getWXACodeUnlimitURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

// QRCoder 小程序码参数
type ResultOfQrCode struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// page 必须是已经发布的小程序存在的页面,根路径前不要填加 /,不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	Page string `json:"page,omitempty"`
	// path 扫码进入的小程序页面路径
	Path string `json:"path,omitempty"`
	// width 图片宽度
	Width int `json:"width,omitempty"`
	// scene 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene,omitempty"`
	// autoColor 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color,omitempty"`
	// lineColor AutoColor 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"},十进制表示
	LineColor Color `json:"line_color,omitempty"`
	// isHyaline 是否需要透明底色
	IsHyaline bool `json:"is_hyaline,omitempty"`
}

func NewQRCode(cfg *config.MiniConfig) *QRCode {
	return &QRCode{
		config: cfg,
	}
}

// fetchCode 请求并返回二维码二进制数据
func (qrCode *QRCode) fetchCode(urlStr string, body interface{}) ([]byte, error) {
	var accessToken string
	tokenInstance := token.NewToken(qrCode.config)
	accessTokenRes, err := tokenInstance.GetToken()
	if err != nil {
		return nil, err
	}
	accessToken = accessTokenRes.AccessToken

	urlStr = fmt.Sprintf(urlStr, accessToken)
	var contentType string
	resp, err := ghttp.Post(urlStr, body)
	if err != nil {
		return nil, err
	}
	response := []byte(resp.RawResponse())
	contentType = resp.Header["Content-Type"][0]
	if contentType == "image/jpeg" {
		// 返回文件
		return response, nil
	}
	result := &ResultOfQrCode{}
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("mini program qrcode fetchCode报文内容解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("mini program qrcode fetchCode报文内容解析失败，error : %v", err))
	}
	return response, nil
}

// CreateWXAQRCode 获取小程序二维码，适用于需要的码数量较少的业务场景
// 文档地址： https://developers.weixin.qq.com/miniprogram/dev/api/createWXAQRCode.html
func (qrCode *QRCode) CreateWXAQRCode(coderParams ResultOfQrCode) (response []byte, err error) {
	return qrCode.fetchCode(createWXAQRCodeURL, coderParams)
}

// GetWXACode 获取小程序码，适用于需要的码数量较少的业务场景
// 文档地址： https://developers.weixin.qq.com/miniprogram/dev/api/getWXACode.html
func (qrCode *QRCode) GetWXACode(coderParams ResultOfQrCode) (response []byte, err error) {
	return qrCode.fetchCode(getWXACodeURL, coderParams)
}

// GetWXACodeUnlimit 获取小程序码，适用于需要的码数量极多的业务场景
// 文档地址： https://developers.weixin.qq.com/miniprogram/dev/api/getWXACodeUnlimit.html
func (qrCode *QRCode) GetWXACodeUnlimit(coderParams ResultOfQrCode) (response []byte, err error) {
	return qrCode.fetchCode(getWXACodeUnlimitURL, coderParams)
}
