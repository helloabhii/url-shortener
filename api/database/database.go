package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8" //here we use this database  //v8 stable version
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client { //this func returns -> redis.Client

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), //created in env file
		Password: os.Getenv("DB_PASS"), //getting this from our file
		DB:       dbNo,
	})
	return rdb

}
