package bytedance

import (
	"github.com/TarrantHightopp/bytedance/cache"
	"github.com/TarrantHightopp/bytedance/miniprogram"
	"github.com/TarrantHightopp/bytedance/miniprogram/config"
)

type ByteDance struct {
	cache cache.Cache
}

func NewByteDance() *ByteDance {
	return &ByteDance{}
}

// SetCache 设置cache
func (by *ByteDance) SetCache(cache cache.Cache) {
	by.cache = cache
}

// GetMiniProgram 获取小程序的实例
func (by *ByteDance) GetMiniProgram(cfg *config.Config) *miniprogram.MiniProgram {
	if cfg.Cache == nil {
		cfg.Cache = by.cache
	}
	return miniprogram.NewMiniProgram(cfg)
}
