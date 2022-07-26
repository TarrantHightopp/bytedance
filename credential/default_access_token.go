package credential

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TarrantHightopp/bytedance/cache"
	"github.com/TarrantHightopp/bytedance/util"
	"sync"
	"time"
)

// doc https://microapp.bytedance.com/docs/zh-CN/mini-app/develop/server/interface-request-credential/get-access-token

const (
	GrantType = "client_credential"
	// AccessTokenURL 获取access_token的接口
	accessTokenURL = "https://developer.toutiao.com/api/apps/v2/token"
	// CacheKeyMiniProgramPrefix 小程序cache key前缀
	CacheKeyMiniProgramPrefix = "/go/bytedance/miniprogram"
)

// DefaultAccessToken 默认AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewDefaultAccessToken new DefaultAccessToken
func NewDefaultAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &DefaultAccessToken{
		appID:           appID,
		appSecret:       appSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError
	Data ResAccessTokenData `json:"data"`
}

type ResAccessTokenData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken 获取access_token,先从cache中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessToken() (accessToken string, err error) {
	// 先从cache中取
	accessTokenCacheKey := fmt.Sprintf("%s/access_token/%s", ak.cacheKeyPrefix, ak.appID)

	var cacheInfo string
	if cacheInfo, err = ak.cache.Get(accessTokenCacheKey); err == nil {
		return cacheInfo, nil
	}

	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()

	// 双检，防止重复从微信服务器获取
	if cacheInfo, err = ak.cache.Get(accessTokenCacheKey); err == nil {
		return cacheInfo, nil
	}

	// cache失效，从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = GetTokenFromServer(accessTokenURL, ak.appID, ak.appSecret)
	if err != nil {
		return
	}

	expires := resAccessToken.Data.ExpiresIn - 1500
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.Data.AccessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	accessToken = resAccessToken.Data.AccessToken
	return
}

// GetAccessTokenRequest 请求获取access token 参数
type GetAccessTokenRequest struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	GrantType string `json:"grant_type"`
}

// GetTokenFromServer 强制从微信服务器获取token
func GetTokenFromServer(url, appid, secret string) (resAccessToken ResAccessToken, err error) {
	getAccessTokenRequest := GetAccessTokenRequest{
		Appid:     appid,
		Secret:    secret,
		GrantType: GrantType,
	}
	var requestBody []byte
	if requestBody, err = json.Marshal(getAccessTokenRequest); err != nil {
		return resAccessToken, err
	}

	fmt.Println(string(requestBody))
	fmt.Println(url)
	var responseBody []byte
	if responseBody, err = util.PostJSON(url, requestBody); err != nil {
		return resAccessToken, err
	}

	if err = json.Unmarshal(responseBody, &resAccessToken); err != nil {
		return resAccessToken, err
	}

	fmt.Println(resAccessToken)
	if resAccessToken.ErrNo != 0 {
		err = errors.New(resAccessToken.ErrTips)
		return
	}
	return
}
