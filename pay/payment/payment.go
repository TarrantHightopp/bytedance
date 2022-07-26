package payment

import (
	"encoding/json"
	"errors"
	"github.com/TarrantHightopp/bytedance/pay/config"
	"github.com/TarrantHightopp/bytedance/util"
)

var PaymentUrl = `https://developer.toutiao.com/api/apps/ecpay/v1/create_order

`

type Order struct {
	AppID        string `json:"app_id"`                  // 小程序APPID
	OutOrderNo   string `json:"out_order_no"`            // 开发者侧的订单号。 只能是数字、大小写字母_-*且在同一个商户号下唯一
	TotalAmount  uint   `json:"total_amount"`            // 支付价格。 单位为[分]
	Subject      string `json:"subject"`                 // 商品描述。 长度限制不超过 128 字节且不超过 42 字符
	Body         string `json:"body"`                    // 商品详情 长度限制不超过 128 字节且不超过 42 字符
	ValidTime    uint   `json:"valid_time"`              // 订单过期时间(秒)。最小5分钟，最大2天，小于5分钟会被置为5分钟，大于2天会被置为2天
	Sign         string `json:"sign"`                    // 签名
	CpExtra      string `json:"cp_extra,omitempty"`      // 开发者自定义字段，回调原样回传。 超过最大长度会被截断
	NotifyUrl    string `json:"notify_url,omitempty"`    // 商户自定义回调地址，必须以 https 开头，支持 443 端口。 指定时，支付成功后抖音会请求该地址通知开发者
	ThirdPartyId string `json:"thirdparty_id,omitempty"` // 第三方平台服务商 id，非服务商模式留空
	StoreUid     string `json:"store_uid,omitempty"`     // 可用此字段指定本单使用的收款商户号（目前为灰度功能，需要联系平台运营添加白名单；未在白名单的小程序，默认使用该小程序下尾号为0的商户号收款）
	DisableMsg   string `json:"disable_msg,omitempty"`   // 是否屏蔽支付完成后推送用户抖音消息，1-屏蔽 0-非屏蔽，默认为0。 特别注意： 若接入POI, 请传1。因为POI订单体系会发消息，所以不用再接收一次担保支付推送消息
	MsgPage      string `json:"msg_page,omitempty"`      // 担保支付消息跳转页
}

type Payment struct {
	*config.Config
}

type ResponsePay struct {
	util.CommonError
	Data ResponsePayData `json:"data"`
}

type ResponsePayData struct {
	OrderId    string `json:"order_id"`
	OrderToken string `json:"order_token"`
}

func (p *Payment) Pay(data *Order) (responsePay ResponsePay, err error) {
	var byteInfo = make([]byte, 0)

	var paramsMap = make(map[string]interface{}, 0)

	if byteInfo, err = json.Marshal(data); err != nil {
		return responsePay, err
	}

	if err = json.Unmarshal(byteInfo, &paramsMap); err != nil {
		return responsePay, err
	}

	data.Sign = p.Sign(paramsMap)

	if byteInfo, err = json.Marshal(data); err != nil {
		return responsePay, err
	}

	var response = make([]byte, 0)
	if response, err = util.PostJSON(PaymentUrl, byteInfo); err != nil {
		return responsePay, err
	}
	if err = json.Unmarshal(response, &responsePay); err != nil {
		return responsePay, err
	}
	if responsePay.ErrNo != 0 {
		return responsePay, errors.New(responsePay.ErrTips)
	}
	return responsePay, nil
}
