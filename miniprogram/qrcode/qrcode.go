package qrcode

import (
	"encoding/json"
	"fmt"
	"github.com/TarrantHightopp/bytedance/miniprogram/context"
	"github.com/TarrantHightopp/bytedance/util"
	"strings"
)

type QRCode struct {
	*context.Context
}

var QRCodeURL = `https://developer.toutiao.com/api/apps/qrcode`

// NewQRCode 实例
func NewQRCode(context *context.Context) *QRCode {
	qrCode := new(QRCode)
	qrCode.Context = context
	return qrCode
}

// Color QRCode color
type Color struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

// QRCoder 小程序码参数
type QRCoder struct {
	// 服务端 API 调用标识，获取方法
	AccessToken string `json:"access_token"`
	// 是打开二维码的字节系 app 名称，默认为今日头条，取值如下表所示
	AppName string `json:"appname,omitempty"`
	// 小程序/小游戏启动参数，小程序则格式为 encode({path}?{query})，小游戏则格式为 JSON 字符串，默认为空
	Path string `json:"path,omitempty"`
	// 二维码宽度，单位 px，最小 280px，最大 1280px，默认为 430px
	Width string `json:"width,omitempty"`
	// 二维码线条颜色，默认为黑色
	LineColor Color `json:"line_color,omitempty"`
	// 二维码背景颜色，默认为白色
	Background Color `json:"background,omitempty"`
	// 是否展示小程序/小游戏 icon，默认不展示
	SetIcon bool `json:"set_icon"`
}

const (
	TOUTIAO     = "toutiao"      // 今日头条
	TOUTIAOLITE = "toutiao_lite" // 今日头条极速版
	DOUYIN      = "douyin"       // 抖音
	DOUYINLITE  = "douyin_lite"  // 抖音极速版
	PIPIXIA     = "pipixia"      // 皮皮虾
	HUOSHAN     = "huoshan"      // 火山小视频
	XIGUA       = "xigua"        // 西瓜视频
)

func (qrCode *QRCode) FetchCode(data QRCoder) (response []byte, err error) {
	if data.AccessToken == "" {
		if data.AccessToken, err = qrCode.GetAccessToken(); err != nil {
			return
		}
	}
	var contentType string
	if response, contentType, err = util.PostJSONWithRespContentType(QRCodeURL, data); err != nil {
		return response, err
	}

	if strings.HasPrefix(contentType, "application/json") {
		// 返回错误信息
		var result util.CommonError
		err = json.Unmarshal(response, &result)
		if err == nil && result.ErrNo != 0 {
			err = fmt.Errorf("fetchCode error : errcode=%v , errmsg=%v", result.ErrNo, result.ErrTips)
			return nil, err
		}
	}
	// 返回文件
	if contentType == "image/jpeg" {
		return response, nil
	}
	return nil, fmt.Errorf("fetchCode error : unknown response content type - %v", contentType)
}
