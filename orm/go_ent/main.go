package main

import (
	"context"
	"github.com/xiaozefeng/go-example/orm/go_ent/db"
	"log"
)

func main() {
	client, err := db.GetClient()
	if err != nil {
		log.Fatalf("failed opening db conneciton , error: %v \n", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v\n", err)
	}
}
