package usecase

import (
	"testing"

	mockSlack "github.com/walnuts1018/wakatime-to-slack-profile/mock/slack"
)

func TestUsecase_SetUserCustomStatus(t *testing.T) {
	c := mockSlack.NewClient()
	c.SetEmojis([]string{
		"gopher",
		"python",
		"TypeScript",
	})

	u := NewUsecase(nil, nil, c, map[string]string{
		"Go": "gopher",
	})

	type args struct {
		language string
	}
	tests := []struct {
		name      string
		u         *Usecase
		args      args
		wantEmoji string
		wantText  string
		wantErr   bool
	}{
		{
			name: "normal",
			u:    u,
			args: args{
				language: "TypeScript",
			},
			wantEmoji: ":TypeScript:",
			wantText:  "TypeScript",
			wantErr:   false,
		},
		{
			name: "override",
			u:    u,
			args: args{
				language: "Go",
			},
			wantEmoji: ":gopher:",
			wantText:  "Go",
			wantErr:   false,
		},
		{
			name: "lowercase",
			u:    u,
			args: args{
				language: "Python",
			},
			wantEmoji: ":python:",
			wantText:  "Python",
			wantErr:   false,
		},
		{
			name: "not found",
			u:    u,
			args: args{
				language: "Rust",
			},
			wantEmoji: ":question:",
			wantText:  "Rust",
			wantErr:   false,
		},
		{
			name: "empty",
			u:    u,
			args: args{
				language: "",
			},
			wantEmoji: ":namakemono:",
			wantText:  "",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.SetUserCustomStatus(tt.args.language); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.SetUserCustomStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotEmoji, gotText := c.GetStatus()
			if gotEmoji != tt.wantEmoji {
				t.Errorf("Usecase.SetUserCustomStatus() gotEmoji = %v, want %v", gotEmoji, tt.wantEmoji)
			}
			if gotText != tt.wantText {
				t.Errorf("Usecase.SetUserCustomStatus() gotText = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}
