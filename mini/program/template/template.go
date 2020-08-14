/*
@Time : 2020/8/14 12:51 下午
@Author : sunmoon
@File : template
@Software: GoLand
*/
package template

import (
	"errors"
	"fmt"
	"github.com/arieslee/gf-wx/mini/config"
	"github.com/arieslee/gf-wx/mini/program/token"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
)

const (
	AddURL         = "https://api.weixin.qq.com/cgi-bin/wxopen/template/add?access_token=%s"
	DeleteURL      = "https://api.weixin.qq.com/wxaapi/newtmpl/deltemplate?access_token=%s"
	GetCategoryURL = "https://api.weixin.qq.com/wxaapi/newtmpl/getcategory?access_token=%s"
	GetKeywordURL  = "https://api.weixin.qq.com/wxaapi/newtmpl/getpubtemplatekeywords?access_token=%s"
	GetTitleURL    = "https://api.weixin.qq.com/wxaapi/newtmpl/getpubtemplatetitles?access_token=%s"
	GetListURL     = "https://api.weixin.qq.com/wxaapi/newtmpl/gettemplate?access_token=%s"
	SendURL        = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
)

type MiniProgramTemplate struct {
	config *config.MiniConfig
}

func NewTemplate(cfg *config.MiniConfig) *MiniProgramTemplate {
	return &MiniProgramTemplate{
		config: cfg,
	}
}

type AddRequest struct {
	Tid       string `json:"tid"`
	KidList   string `json:"kid_list"` // 如，1,2
	SceneDesc string `json:"scene_desc"`
}
type AddResponse struct {
	ErrCode   int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PriTmplId string `json:"priTmplId"`
}

//
// @See : https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.addTemplate.html
func (t *MiniProgramTemplate) Add(req *AddRequest) (*AddResponse, error) {
	kidList := gstr.SplitAndTrim(req.KidList, ",")
	count := len(kidList)
	if count < 2 || count > 5 {
		return nil, errors.New("最多支持5个，最少2个关键词组合")
	}
	if len(req.SceneDesc) == 0 {
		return nil, errors.New("服务场景描述不能为空")
	}
	length := len([]rune(req.SceneDesc))
	if length > 15 {
		return nil, errors.New("服务场景描述的长度必须15个字以内")
	}
	if len(req.Tid) <= 0 {
		return nil, errors.New("模板标题 id 为能为空")
	}
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(AddURL, tokenData.AccessToken)
	response := ghttp.PostBytes(url, req)
	result := &AddResponse{}
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("AddTemplate报文解析失败，error : %v", err)
		return nil, fmt.Errorf("AddTemplate报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("AddTemplate报文解析失败，error : %v", err)
		return nil, fmt.Errorf("AddTemplate报文解析失败，error : %v", err)
	}
	return result, nil
}

type DeleteResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// @see : https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.deleteTemplate.html
func (t *MiniProgramTemplate) Delete(priTmplId string) (*DeleteResponse, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(DeleteURL, tokenData.AccessToken)
	result := &DeleteResponse{}
	response := ghttp.PostBytes(url, g.Map{
		"priTmplId": priTmplId,
	})
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("DeleteTemplate报文解析失败，error : %v", err)
		return nil, fmt.Errorf("DeleteTemplate报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("DeleteTemplate报文解析失败，error : %v", err)
		return nil, fmt.Errorf("DeleteTemplate报文解析失败，error : %v", err)
	}
	return result, nil
}

type CategoryItem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type CategoryList struct {
	ErrCode int64           `json:"errcode"`
	ErrMsg  string          `json:"errmsg"`
	Data    []*CategoryItem `json:"data"`
}

// @see: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getCategory.html
func (t *MiniProgramTemplate) GetCategory() (*CategoryList, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(GetCategoryURL, tokenData.AccessToken)
	result := &CategoryList{}
	response := ghttp.GetBytes(url)
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetCategory报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetCategory报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("GetCategory报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetCategory报文解析失败，error : %v", err)
	}
	return result, nil
}

type KeywordItem struct {
	Kid     uint   `json:"kid"`
	Name    string `json:"name"`
	Example string `json:"example"`
	Rule    string `json:"rule"`
}

type Keywords struct {
	ErrCode int64          `json:"errcode"`
	ErrMsg  string         `json:"errmsg"`
	Count   uint           `json:"count"`
	Data    []*KeywordItem `json:"data"`
}

// @see : https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getPubTemplateKeyWordsById.html
func (t *MiniProgramTemplate) GetKeywordById(tid string) (*Keywords, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(GetKeywordURL, tokenData.AccessToken)
	result := &Keywords{}
	response := ghttp.GetBytes(url, g.Map{
		"tid": tid,
	})
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetKeywordById报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetKeywordById报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("GetKeywordById报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetKeywordById报文解析失败，error : %v", err)
	}
	return result, nil
}

type TitleItem struct {
	Tid        int    `json:"tid"`
	Title      string `json:"title"`
	Type       uint   `json:"type"`
	CategoryId int    `json:"categoryId"`
}

type TitleList struct {
	ErrCode int64        `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
	Count   uint         `json:"count"`
	Data    []*TitleItem `json:"data"`
}
type TitleRequest struct {
	Ids   string `json:"ids"`
	Start uint   `json:"start"`
	Limit uint   `json:"limit"`
}

// @See: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getPubTemplateTitleList.html
func (t *MiniProgramTemplate) GetTitleList(req *TitleRequest) (*TitleList, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(GetKeywordURL, tokenData.AccessToken)
	result := &TitleList{}
	response := ghttp.GetBytes(url, req)
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetTitleList报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetTitleList报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("GetTitleList报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetTitleList报文解析失败，error : %v", err)
	}
	return result, nil
}

type TplItem struct {
	PriTmplId string `json:"priTmplId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Example   string `json:"example"`
	Type      uint   `json:"type"`
}
type TplList struct {
	ErrCode int64      `json:"errcode"`
	ErrMsg  string     `json:"errmsg"`
	Data    []*TplItem `json:"data"`
}

// @see: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getTemplateList.html
func (t *MiniProgramTemplate) GetTemplateList() (*TplList, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(GetListURL, tokenData.AccessToken)
	result := &TplList{}
	response := ghttp.GetBytes(url)
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("GetTemplateList报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetTemplateList报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("GetTemplateList报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetTemplateList报文解析失败，error : %v", err)
	}
	return result, nil
}

// @see : https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
type SendRequest struct {
	ToUser           string          `json:"touser"`
	TemplateId       string          `json:"template_id"`
	Page             string          `json:"page"`
	MiniProgramState string          `json:"miniprogram_state"`
	Lang             string          `json:"lang"`
	Data             []*SendDataItem `json:"data"` //模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }
}

// "key1": { "value": any }
type SendDataItem struct {
	Key *SendDataItemValue `json:"key"`
}

// { "value": any }
type SendDataItemValue struct {
	Value string `json:"value"`
}
type SendResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// @see: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
func (t *MiniProgramTemplate) Send(req *SendRequest) (*SendResponse, error) {
	newToken := token.NewToken(t.config)
	tokenData, err := newToken.GetToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(GetListURL, tokenData.AccessToken)
	result := &SendResponse{}
	response := ghttp.GetBytes(url)
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Fatalf("SendTemplate报文解析失败报文解析失败，error : %v", err)
		return nil, fmt.Errorf("SendTemplate报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Fatalf("SendTemplate报文解析失败，error : %v", err)
		return nil, fmt.Errorf("SendTemplate报文解析失败，error : %v", err)
	}
	return result, nil
}
