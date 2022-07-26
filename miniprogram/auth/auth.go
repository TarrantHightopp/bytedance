package auth

import (
	"encoding/json"
	"errors"
	"github.com/TarrantHightopp/bytedance/miniprogram/context"
	"github.com/TarrantHightopp/bytedance/util"
)

const (
	Code2SessionUrl string = `https://developer.toutiao.com/api/apps/v2/jscode2session`
)

// Auth 登录/用户信息
type Auth struct {
	*context.Context
}

// doc https://microapp.bytedance.com/docs/zh-CN/mini-app/develop/server/log-in/code-2-session

// RequestCode2Session 登录凭证校验的请求参数
type RequestCode2Session struct {
	Appid         string `json:"appid"`
	Secret        string `json:"secret"`
	AnonymousCode string `json:"anonymous_code"`
	Code          string `json:"code"`
}

// ResponseCode2Session 登录凭证校验的返回结果
type ResponseCode2Session struct {
	util.CommonError
	Data ResponseCode2SessionData `json:"data"`
}

type ResponseCode2SessionData struct {
	SessionKey      string `json:"session_key"`
	Openid          string `json:"openid"`
	AnonymousOpenid string `json:"anonymous_openid"`
	Unionid         string `json:"unionid"`
}

// NewAuth new auth
func NewAuth(ctx *context.Context) *Auth {
	return &Auth{ctx}
}

// Code2SessionContext 登录。
func (auth *Auth) Code2SessionContext(anonymousCode, code string) (result ResponseCode2Session, err error) {
	requestCode2Session := RequestCode2Session{
		Appid:         auth.AppID,
		Secret:        auth.AppSecret,
		AnonymousCode: anonymousCode,
		Code:          code,
	}
	var requestBody []byte
	if requestBody, err = json.Marshal(requestCode2Session); err != nil {
		return result, err
	}

	var responseBody []byte
	if responseBody, err = util.PostJSON(Code2SessionUrl, requestBody); err != nil {
		return result, err
	}

	if err = json.Unmarshal(responseBody, &result); err != nil {
		return result, err
	}

	if result.ErrNo != 0 {
		return result, errors.New(result.ErrTips)
	}
	return
}
