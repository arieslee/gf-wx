/*
@Time : 2020/7/31 3:38 下午
@Author : sunmoon
@File : js
@Software: GoLand
*/
package js

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mp/config"
	"github.com/arieslee/gf-wx/mp/token"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"sync"
)

type Js struct {
	config *config.MpConfig
}

var jsAPITicketSync *sync.RWMutex

const (
	jsTicketCacheKey = "gf-wx-js-ticket:%s"
	getTicketURL     = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// TicketResult 请求jsapi_tikcet返回结果
type TicketResult struct {
	ErrCode   int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

// 生成的config
type Config struct {
	AppID     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonce_str"`
	Signature string `json:"signature"`
}

func NewJS(cfg *config.MpConfig) *Js {
	InitSync()
	return &Js{
		config: cfg,
	}
}

func InitSync() {
	jsAPITicketSync = new(sync.RWMutex)
}

func (j *Js) GetConfig(url string) (*Config, error) {
	ticketResult, err := j.Ticket()
	if err != nil {
		return nil, err
	}
	jsTicket := ticketResult.Ticket
	jsConfig := &Config{}
	nonce := grand.S(32)
	timestamp := gtime.Now().Unix()
	sign := j.GetSignature(jsTicket, nonce, timestamp, url)
	jsConfig.Signature = sign
	jsConfig.AppID = j.config.AppID
	jsConfig.Timestamp = timestamp
	jsConfig.NonceStr = nonce
	return jsConfig, nil
}

// 获取json格式的config
func (j *Js) BuildConfigStr(url string) string {
	cfg, _ := j.GetConfig(url)
	str, _ := gjson.Encode(cfg)
	return gconv.String(str)
}

func (j *Js) GetSignature(ticket, nonce string, timestamp int64, url string) string {
	return fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, url)
}

func (j *Js) Ticket() (*TicketResult, error) {
	jsAPITicketSync.Lock()
	defer jsAPITicketSync.Unlock()
	key := fmt.Sprintf(jsTicketCacheKey, j.config.AppID)
	cacheData, _ := g.Redis().Do("GET", key)
	ticketStr := gconv.String(cacheData)
	if len(ticketStr) <= 0 {
		tokenStruct := token.NewToken(j.config)
		tokenResult, err := tokenStruct.GetToken()
		if err != nil {
			fmt.Println(err.Error())
		}
		accessToken := tokenResult.AccessToken
		return j.GetTicketFromServer(accessToken)
	}
	result := &TicketResult{}
	err := gjson.DecodeTo(ticketStr, &result)
	if err != nil {
		glog.Line().Fatalf("Ticket 缓存内容解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("Ticket 缓存内容解析失败，error : %v", err))
	}
	return result, nil
}
func (j *Js) GetTicketFromServer(accessToken string) (*TicketResult, error) {
	url := fmt.Sprintf(getTicketURL, accessToken)
	response := ghttp.GetBytes(url)
	result := &TicketResult{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetTicketFromServer报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("GetTicketFromServer报文解析失败，error : %v", err))
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("GetTicketFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, errors.New(fmt.Sprintf("GetTicketFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
	}
	value := gconv.Map(result)
	expire := result.ExpiresIn - 100
	key := fmt.Sprintf(jsTicketCacheKey, j.config.AppID)
	g.Redis().Do("SETEX", key, expire, gconv.String(value))
	return result, nil
}
