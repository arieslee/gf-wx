/*
@Time : 2020/8/12 1:02 下午
@Author : sunmoon
@File : config
@Software: GoLand
*/
package config

type MiniConfig struct {
	AppID          string `json:"app_id"`           //appid
	AppSecret      string `json:"app_secret"`       //appsecret
	Token          string `json:"token"`            //token
	MchID          string `json:"mch_id"`           // 商户号
	ParamSignKey   string `json:"param_sign_key"`   //  参数加密密钥
	EncodingAESKey string `json:"encoding_aes_key"` //EncodingAESKey
}
