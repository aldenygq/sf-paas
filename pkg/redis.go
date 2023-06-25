package pkg

import (
	"fmt"
	"sf-paas/config"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	ip := config.Conf.Redis.Ip
	port := config.Conf.Redis.Port
	log.Info("ip:", ip)
	log.Info("port:", port)
	url := fmt.Sprintf("%s:%s", ip, port)
	r := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	RedisClient = r
	if RedisClient == nil {
		log.Error("init redis client fialed")
		return
	}
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Error("connect redis failed:", err)
		return
	}
}
