/*
@Time : 2020/7/30 6:09 下午
@Author : sunmoon
@File : config
@Software: GoLand
*/
package config

type MpConfig struct {
    AppID          string `json:"app_id"`           //appid
    AppSecret      string `json:"app_secret"`       //appsecret
    Token          string `json:"token"`            //token
    EncodingAESKey string `json:"encoding_aes_key"` //EncodingAESKey
}