package miniprogram

import (
	"github.com/TarrantHightopp/bytedance/credential"
	"github.com/TarrantHightopp/bytedance/miniprogram/auth"
	"github.com/TarrantHightopp/bytedance/miniprogram/config"
	"github.com/TarrantHightopp/bytedance/miniprogram/context"
)

type MiniProgram struct {
	ctx *context.Context
}

func NewMiniProgram(cfg *config.Config) *MiniProgram {
	defaultAkHandle := credential.NewDefaultAccessToken(cfg.AppID, cfg.AppSecret, credential.CacheKeyMiniProgramPrefix, cfg.Cache)
	ctx := &context.Context{
		Config:            cfg,
		AccessTokenHandle: defaultAkHandle,
	}
	return &MiniProgram{ctx}
}

func (miniProgram *MiniProgram) GetAuth() *auth.Auth {
	return auth.NewAuth(miniProgram.ctx)
}
