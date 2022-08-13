package user

import (
	"context"
	"fmt"
	"github.com/xiaozefeng/go-example/orm/go_ent/ent"
	"log"
)

func CreateUser(ctx context.Context, user ent.User, client *ent.Client) (*ent.User, error) {

	u, err := client.User.
		Create().
		SetAge(user.Age).
		SetName(user.Name).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created", u)
	return u, nil
}
