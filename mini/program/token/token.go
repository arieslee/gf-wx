/*
@Time : 2020/7/31 3:44 下午
@Author : sunmoon
@File : accessToken
@Software: GoLand
*/
package token

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"sync"
)

type MiniProgramToken struct {
	config *config.MiniConfig
}

const (
	getTokenURL   = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	tokenCacheKey = "gf-wx-token:%s"
)

func NewToken(cfg *config.MiniConfig) *MiniProgramToken {
	return &MiniProgramToken{
		config: cfg,
	}
}

//ResAccessToken struct
type AccessTokenResult struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

var syncLock *sync.Mutex

func InitSyncLock() {
	syncLock = new(sync.Mutex)
}
func (t *MiniProgramToken) GetToken() (*AccessTokenResult, error) {
	key := fmt.Sprintf(tokenCacheKey, t.config.AppID)
	cacheData, _ := g.Redis().Do("GET", key)
	tokenStr := gconv.String(cacheData)
	result := &AccessTokenResult{}
	if len(tokenStr) <= 0 {
		return t.GetTokenFromServer()
	}
	err := gjson.DecodeTo(tokenStr, &result)
	if err != nil {
		g.Log().Line().Fatalf("GetToken缓存内容解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("GetToken缓存内容解析失败，error : %v", err))
	}
	return result, nil
}

func (t *MiniProgramToken) GetTokenFromServer() (*AccessTokenResult, error) {
	if syncLock != nil {
		syncLock.Lock()
		defer syncLock.Unlock()
	}
	url := fmt.Sprintf(getTokenURL, t.config.AppID, t.config.AppSecret)
	response := ghttp.GetBytes(url)
	result := &AccessTokenResult{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		g.Log().Line().Fatalf("GetTokenFromServer报文解析失败，error : %v", err)
		return nil, errors.New(fmt.Sprintf("GetTokenFromServer报文解析失败，error : %v", err))
	}
	if result.ErrCode != 0 {
		g.Log().Line().Fatalf("GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, errors.New(fmt.Sprintf("GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg))
	}
	value := gconv.Map(result)
	expire := result.ExpiresIn - 100
	key := fmt.Sprintf(tokenCacheKey, t.config.AppID)
	g.Redis().Do("SETEX", key, expire, gconv.String(value))
	return result, nil
}
