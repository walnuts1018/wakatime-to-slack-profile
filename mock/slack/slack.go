package mockSlack

import (
	"fmt"
	"strings"
)

type client struct {
	emojis []string
	emoji  string
	text   string
}

func NewClient() *client {
	return &client{}
}

func (c *client) SetEmojis(emojis []string) {
	emojis = append(emojis, "namakemono", "question")
	c.emojis = emojis
}

func (c *client) SetUserCustomStatus(emoji string, text string) error {
	for _, e := range c.emojis {
		if e == emoji {
			if !(strings.HasPrefix(emoji, ":") && strings.HasSuffix(emoji, ":")) {
				emoji = ":" + emoji + ":"
			}

			c.emoji = emoji
			c.text = text
			return nil
		}
	}
	return fmt.Errorf("no such emoji: %v", emoji)
}

func (c *client) GetStatus() (string, string) {
	return c.emoji, c.text
}
