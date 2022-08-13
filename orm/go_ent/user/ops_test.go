package user

import (
	"context"
	"github.com/xiaozefeng/go-example/orm/go_ent/db"
	"github.com/xiaozefeng/go-example/orm/go_ent/ent"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client, err := db.GetClient()
	if err != nil {
		t.Error(err)
	}
	u := ent.User{
		Age:  20,
		Name: "micket",
	}
	user, err := CreateUser(context.Background(), u, client)
	if err != nil {
		t.Error(err)
	}
	if user.Name != u.Name {
		t.Errorf("expect: %s, but got: %s", u.Name, user.Name)
	}
	if u.Age != user.Age {
		t.Errorf("Expect: %d but go: %d", u.Age, user.Age)
	}
}
