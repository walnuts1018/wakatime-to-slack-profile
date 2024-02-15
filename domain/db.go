package domain

type DB interface {
	SaveOAuth2Token(OAuth2Token) error
	GetOAuth2Token() (OAuth2Token, error)
	UpdateOAuth2Token(token OAuth2Token) error
	Close() error
}
