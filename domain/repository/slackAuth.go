package repository

type SlackAuth interface {
	// GetOAuthV2Response returns the access token and the user ID
	GetOAuthV2Response(code string) (string, string, error)
}
