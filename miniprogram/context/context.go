package context

import (
	"github.com/TarrantHightopp/bytedance/credential"
	"github.com/TarrantHightopp/bytedance/miniprogram/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
