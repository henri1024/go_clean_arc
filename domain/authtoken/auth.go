package authtoken

type Token struct {
	AccessToken    string
	RefreshToken   string
	TokenUuid      string
	RefreshUuid    string
	AccessExpired  int64
	RefreshExpired int64
}
