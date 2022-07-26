package credential

type AccessTokenHandle interface {
	GetAccessToken() (accessToken string, err error)
}
