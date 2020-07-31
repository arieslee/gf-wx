/*
@Time : 2020/7/31 4:50 下午
@Author : sunmoon
@File : server
@Software: GoLand
*/
package server

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"io"
	"reflect"
	"sort"
)

type Server struct {
	Response *ghttp.Response
	Request  *ghttp.Request
	Debug    bool
	Token    string
}

func NewServer(resp *ghttp.Response, req *ghttp.Request, token string) *Server {
	return &Server{
		Response: resp,
		Request:  req,
		Debug:    false,
		Token:    token,
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

}
func (s *Server) requestHandler() {

}
func (s *Server) buildResponse(reply *Reply) error {
	msgType := reply.MsgType
	switch msgType {
	case MsgTypeText:
	case MsgTypeImage:
	case MsgTypeVoice:
	case MsgTypeVideo:
	case MsgTypeMusic:
	case MsgTypeNews:
	default:
		return ErrUnsupportReply
	}

	msgData := reply.MsgData
	value := reflect.ValueOf(msgData)
	//msgData must be a ptr
	kind := value.Kind().String()
	if "ptr" != kind {
		return ErrUnsupportReply
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(s.requestMsg.FromUserName)
	value.MethodByName("SetToUserName").Call(params)

	params[0] = reflect.ValueOf(s.requestMsg.ToUserName)
	value.MethodByName("SetFromUserName").Call(params)

	params[0] = reflect.ValueOf(msgType)
	value.MethodByName("SetMsgType").Call(params)

	params[0] = reflect.ValueOf(util.GetCurrTs())
	value.MethodByName("SetCreateTime").Call(params)

	srv.responseMsg = msgData
	srv.responseRawXMLMsg, err = xml.Marshal(msgData)
	return
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
