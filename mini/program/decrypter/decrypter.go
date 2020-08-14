/*
@Time : 2020/8/12 1:19 下午
@Author : sunmoon
@File : crypto
@Software: GoLand
*/
package decrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/arieslee/gf-wx/mini/config"
)

type MiniProgramDecrypter struct {
	config *config.MiniConfig
}

func NewDecrypter(cfg *config.MiniConfig) *MiniProgramDecrypter {
	return &MiniProgramDecrypter{
		config: cfg,
	}
}

var (
	// ErrAppIDNotMatch appid不匹配
	ErrAppIDNotMatch = errors.New("appID 不匹配")
	// ErrInvalidBlockSize block size不合法
	ErrInvalidBlockSize = errors.New("无效的块大小")
	// ErrInvalidPKCS7Data PKCS7数据不合法
	ErrInvalidPKCS7Data = errors.New("无效的PKCS7数据")
	// ErrInvalidPKCS7Padding 输入padding失败
	ErrInvalidPKCS7Padding = errors.New("输入的无效填充")
)

// PlainResult 解密后的原始结果
type PlainResult struct {
	OpenID          string `json:"openId"`
	UnionID         string `json:"unionId"`
	NickName        string `json:"nickName"`
	Gender          int    `json:"gender"`
	City            string `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	AvatarURL       string `json:"avatarUrl"`
	Language        string `json:"language"`
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

// @from https://github.com/silenceper/wechat/blob/release-2.0/miniprogram/encryptor/encryptor.go
// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

// @from https://github.com/silenceper/wechat/blob/release-2.0/miniprogram/encryptor/encryptor.go
// getCipherText returns slice of the cipher text
func getCipherText(sessionKey, encryptedData, iv string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// Decrypt 解密数据
func (decrypter *MiniProgramDecrypter) Decrypt(sessionKey, encryptedData, iv string) (*PlainResult, error) {
	cipherText, err := getCipherText(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}
	var plainResult PlainResult
	err = json.Unmarshal(cipherText, &plainResult)
	if err != nil {
		return nil, err
	}
	if plainResult.Watermark.AppID != decrypter.config.AppID {
		return nil, ErrAppIDNotMatch
	}
	return &plainResult, nil
}
