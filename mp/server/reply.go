/*
@Time : 2020/7/31 5:07 下午
@Author : sunmoon
@File : echo
@Software: GoLand
*/
package server

import (
	"encoding/xml"
	"errors"
)

// MsgType 消息类型
type MsgType string
type EventType string
type CDATA string

var (
	ErrInvalidReply   = errors.New("无效的回复消息")
	ErrUnsupportReply = errors.New("不支持的回复消息")
)

// 微信支持的消息类型
const (
	MsgTypeText       MsgType = "text"       // 文本消息
	MsgTypeImage      MsgType = "image"      // 图片消息
	MsgTypeVoice      MsgType = "voice"      // 语音消息
	MsgTypeVideo      MsgType = "video"      // 视频消息
	MsgTypeShortVideo MsgType = "shortvideo" // 小视频消息
	MsgTypeLocation   MsgType = "location"   // 地理位置消息
	MsgTypeLink       MsgType = "link"       // 链接消息
	MsgTypeMusic      MsgType = "music"      // 音乐消息
	MsgTypeNews       MsgType = "news"       // 图文消息
	MsgTypeWXCard     MsgType = "wxcard"     // 卡券，客服消息时使用
	MsgTypeEvent      MsgType = "event"      // 事件推送
)

// 事件
const (
	//EventSubscribe 订阅
	EventSubscribe = "subscribe"
	//EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	//EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	//EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	//EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	//EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	//EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	//EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	//EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	//EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	//EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	//EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
	//EventTemplateSendJobFinish 发送模板消息推送通知
	EventTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
	//EventWxaMediaCheck 异步校验图片/音频是否含有违法违规内容推送事件
	EventWxaMediaCheck = "wxa_media_check"
)

// 加密回复体
type EncryptMessage struct {
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    string
	Nonce        CDATA
}

// 基本内容
type BaseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      MsgType  `xml:"MsgType"`
}

// EventPic 发图事件推送
type EventPic struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

// 事件消息
type EventMessage struct {
	Event       EventType `xml:"Event"`
	EventKey    string    `xml:"EventKey"`
	Ticket      string    `xml:"Ticket"`
	Latitude    string    `xml:"Latitude"`
	Longitude   string    `xml:"Longitude"`
	Precision   string    `xml:"Precision"`
	MenuID      string    `xml:"MenuId"`
	Status      string    `xml:"Status"`
	SessionFrom string    `xml:"SessionFrom"`

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"`
		ScanResult string `xml:"ScanResult"`
	} `xml:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int32      `xml:"Count"`
		PicList []EventPic `xml:"PicList>item"`
	} `xml:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X"`
		LocationY float64 `xml:"Location_Y"`
		Scale     float64 `xml:"Scale"`
		Label     string  `xml:"Label"`
		Poiname   string  `xml:"Poiname"`
	}
}

// 考卷消息
type CardMessage struct {
	CardID              string `xml:"CardId"`
	RefuseReason        string `xml:"RefuseReason"`
	IsGiveByFriend      int32  `xml:"IsGiveByFriend"`
	FriendUserName      string `xml:"FriendUserName"`
	UserCardCode        string `xml:"UserCardCode"`
	OldUserCardCode     string `xml:"OldUserCardCode"`
	OuterStr            string `xml:"OuterStr"`
	IsRestoreMemberCard int32  `xml:"IsRestoreMemberCard"`
	UnionID             string `xml:"UnionId"`
}

// DeviceMessage 设备消息响应
type DeviceMessage struct {
	DeviceType string
	DeviceID   string
	SessionID  string
	OpenID     string
}

// AllMessage 存放所有微信发送过来的消息和事件
type AllMessage struct {
	BaseMessage
	//基本消息
	MsgID        int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	Recognition  string  `xml:"Recognition"`
	PicURL       string  `xml:"PicUrl"`
	MediaID      string  `xml:"MediaId"`
	Format       string  `xml:"Format"`
	ThumbMediaID string  `xml:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64 `xml:"Scale"`
	Label        string  `xml:"Label"`
	Title        string  `xml:"Title"`
	Description  string  `xml:"Description"`
	URL          string  `xml:"Url"`
	//事件相关
	EventMessage
	// 卡券相关
	CardMessage
	//设备相关
	DeviceMessage
}
