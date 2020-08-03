/*
@Time : 2020/7/31 4:50 下午
@Author : sunmoon
@File : server
@Software: GoLand
*/
package server

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"io"
	"log"
	"sort"
	"strings"
)

type Server struct {
	Response       *ghttp.Response
	Request        *ghttp.Request
	Debug          bool
	AppID          string
	Token          string
	EncodingAESKey string
}

func NewServer(resp *ghttp.Response, req *ghttp.Request, appID, token, aesKey string) *Server {
	return &Server{
		Response:       resp,
		Request:        req,
		Debug:          false,
		AppID:          appID,
		Token:          token,
		EncodingAESKey: aesKey,
	}
}

// 监听
func (s Server) Monitor() error {
	if !s.Validate() {
		return fmt.Errorf("请求校验失败")
	}
	echoStr := s.Request.GetString("echostr", "")
	if len(echoStr) > 0 {
		s.Response.WriteHeader(200)
		s.Response.Write(echoStr)
		return nil
	}
	switch encryptType := s.Request.GetString("encrypt_type"); encryptType {
	case "aes":
		if len(s.EncodingAESKey) <= 0 {
			log.Println("EncodingAESKey无效")
		}
	}
	return nil
}

// DecryptMsg 解密微信消息,密文string->base64Dec->aesDec->去除头部随机字串
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成
func (s *Server) DecryptMessage(msg string) (string, error) {
	aesMsg, err := gbase64.DecodeString(msg)
	if err != nil {
		return "", err
	}
	buf, err := gaes.Decrypt(aesMsg, gconv.Bytes(s.EncodingAESKey))
	var msgLen int32
	binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &msgLen)
	if msgLen < 0 || msgLen > 1000000 {
		return "", errors.New("AesKey is invalid")
	}
	if string(buf[20+msgLen:]) != s.AppID {
		return "", errors.New("AppId is invalid")
	}
	return string(buf[20 : 20+msgLen]), nil
}

// CDATA 标准规范，XML编码成 `<![CDATA[消息内容]]>`
type CDATA string

// wxRespEnc 加密回复体
type wxRespEnc struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    string
	Nonce        CDATA
}

// EncryptMsg 加密普通回复(AES-CBC),打包成xml格式
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppId组成
func (s *Server) EncryptMsg(msg []byte, timeStamp, nonce string) (re *wxRespEnc, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return
	}
	l := buf.Bytes()

	rd := []byte(grand.S(16))

	plain := bytes.Join([][]byte{rd, l, msg, []byte(s.AppID)}, nil)
	ae, _ := gaes.Decrypt(plain, gconv.Bytes(s.EncodingAESKey))
	encMsg := base64.StdEncoding.EncodeToString(ae)
	re = &wxRespEnc{
		Encrypt:      CDATA(encMsg),
		MsgSignature: CDATA(s.makeSignature(s.Token, timeStamp, nonce, encMsg)),
		TimeStamp:    timeStamp,
		Nonce:        CDATA(nonce),
	}
	return
}
func (s *Server) makeSignature(str ...string) string {
	sort.Strings(str)
	h := sha1.New()
	h.Write([]byte(strings.Join(str, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SetDebug set debug field
func (s *Server) SetDebug(debug bool) {
	s.Debug = debug
}

//Validate 校验请求是否合法
func (s *Server) Validate() bool {
	if s.Debug {
		return true
	}
	timestamp := s.Request.GetString("timestamp")
	nonce := s.Request.GetString("nonce")
	signature := s.Request.GetString("signature")
	return signature == s.Signature(s.Token, timestamp, nonce)
}

//Signature sha1签名
func (s *Server) Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		_, _ = io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
