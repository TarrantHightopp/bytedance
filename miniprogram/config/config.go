package config

import "github.com/TarrantHightopp/bytedance/cache"

type Config struct {
	AppID     string `json:"app_id"`     // 小程序 ID
	AppSecret string `json:"app_secret"` // 小程序的 APP Secret 可以在开发者后台获取
	Cache     cache.Cache
}
