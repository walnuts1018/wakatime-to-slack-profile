package domain

type SlackClient interface {
	SetUserCustomStatus(emoji string) error
}
