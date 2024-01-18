package domain

type SlackClient interface {
	SetUserCustomStatus(emoji string, text string) error
}
