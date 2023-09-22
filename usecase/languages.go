package usecase

import (
	"context"
	"fmt"
)

func (u *Usecase) Languages(ctx context.Context) error {
	langs, err := u.wakatimeClient.Languages(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", langs)
	return nil
}
