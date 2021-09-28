package redisDriver

import (
	"context"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type RedisDriver struct {
	rdb *redis.Client
	username string
	password string
}

func (d *RedisDriver) Init(name , password string)  {
	d.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	d.password = password
	d.username = name
}

func (d *RedisDriver) GetUser(ctx context.Context) string{
	rdb := d.rdb
	val , err := rdb.Get(ctx, d.username).Result()
	if err != nil {
		return ""
	}
	return val
}

func (d *RedisDriver) Del(ctx context.Context,string string) {
	_ = d.rdb.Del(ctx, string)
}


func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (d *RedisDriver)CheckPasswordHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(d.password))
	return err == nil
}

func (d *RedisDriver) PutUser(ctx context.Context)  bool{
	if d.GetUser(ctx) == ""{
		p,er := hashPassword(d.password)
		if er!= nil {
			log.Fatal(er)
			return false
		}
		err := d.rdb.Set(ctx, d.username, p, 600*time.Second).Err()
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	return false
}

