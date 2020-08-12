/*
@Time : 2020/8/12 1:59 下午
@Author : sunmoon
@File : payment
@Software: GoLand
*/
package payment

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/helper"
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"strconv"
)

type Payment struct {
	config *config.MiniConfig
}

const (
	PrePareURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

func NewPayment(cfg *config.MiniConfig) *Payment {
	return &Payment{
		config: cfg,
	}
}

// UnifiedOrderResponse 是 unifiedorder 接口的返回
type UnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	TradeType  string `xml:"trade_type,omitempty"`
	PrePayID   string `xml:"prepay_id,omitempty"`
	CodeURL    string `xml:"code_url,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

// PaymentRequest 接口请求参数
type PayRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	DeviceInfo     string `xml:"device_info,omitempty"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	SignType       string `xml:"sign_type,omitempty"`
	Body           string `xml:"body"`
	Detail         string `xml:"detail,omitempty"`
	Attach         string `xml:"attach,omitempty"`      // 附加数据
	OutTradeNo     string `xml:"out_trade_no"`          // 商户订单号
	FeeType        string `xml:"fee_type,omitempty"`    // 标价币种
	TotalFee       string `xml:"total_fee"`             // 标价金额
	SpbillCreateIP string `xml:"spbill_create_ip"`      // 终端IP
	TimeStart      string `xml:"time_start,omitempty"`  // 交易起始时间
	TimeExpire     string `xml:"time_expire,omitempty"` // 交易结束时间
	GoodsTag       string `xml:"goods_tag,omitempty"`   // 订单优惠标记
	NotifyURL      string `xml:"notify_url"`            // 通知地址
	TradeType      string `xml:"trade_type"`            // 交易类型
	ProductID      string `xml:"product_id,omitempty"`  // 商品ID
	LimitPay       string `xml:"limit_pay,omitempty"`   //
	OpenID         string `xml:"openid,omitempty"`      // 用户标识
	SceneInfo      string `xml:"scene_info,omitempty"`  // 场景信息

	XMLName struct{} `xml:"xml"`
}

func (p *Payment) GenPrePareOrder(params *PayRequest) (*UnifiedOrderResponse, error) {
	if len(params.OutTradeNo) <= 0 {
		return nil, errors.New("缺少统一支付接口必填参数out_trade_no！")
	}
	if len(params.Body) <= 0 {
		return nil, errors.New("缺少统一支付接口必填参数body！")
	}
	if len(params.TotalFee) <= 0 {
		return nil, errors.New("缺少统一支付接口必填参数total_fee！")
	}
	if len(params.TradeType) <= 0 {
		return nil, errors.New("缺少统一支付接口必填参数trade_type！")
	}
	if len(params.OpenID) <= 0 {
		return nil, errors.New("统一支付接口中，缺少必填参数openid！")
	}
	if len(params.NotifyURL) <= 0 {
		return nil, errors.New("统一支付接口中，缺少必填参数NotifyURL！")
	}
	// 自动把元转成分
	totalFee := gconv.Float64(params.TotalFee) * 100
	params.TotalFee = strconv.FormatInt(gconv.Int64(totalFee), 10)
	params.NonceStr = grand.S(32)
	params.AppID = p.config.AppID
	params.MchID = p.config.MchID
	params.SignType = helper.SignTypeMD5
	ip, err := gipv4.GetIpArray()
	if err != nil {
		return nil, fmt.Errorf("无法获取服务器的IP,err:%v", err)
	}
	params.SpbillCreateIP = ip[0]
	postMap := gconv.MapStrStr(params)
	sign, err := helper.ParamSign(postMap, p.config.ParamSignKey)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名时发生错误,err:%v", err)
	}
	params.Sign = sign
	toXmlData := gconv.Map(params)
	xmlData, _ := gxml.Encode(toXmlData)
	response := ghttp.PostBytes(PrePareURL, gconv.String(xmlData))
	result := &UnifiedOrderResponse{}
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ReturnCode != "SUCCESS" {
		glog.Line().Fatalf("支付失败，error : %s", result.ReturnMsg)
		return nil, fmt.Errorf("支付失败，error : %s", result.ReturnMsg)
	}
	return result, nil
}
