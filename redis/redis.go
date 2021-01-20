package redis

import (
	"fmt"
	"time"
	"usermana/config"

	"github.com/go-redis/redis"
)

var client *redis.Client

//init  redis的初始化函数.
func init() {
	//连接redis服务器.
	client = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		// 数据库
		DB: 0,
		// 连接池大小 Maximum number of socket connections.
		PoolSize: config.RedisPoolSize,
	})
	//测试redis连接是否成功.
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis init done.")
}

// GetProfile 获取用户信息.
func GetProfile(userName string) (nickName string, picName string, hasData bool, err error) {
	vals, err := client.HGetAll(client.Context(), userName).Result()
	if err != nil {
		return "", "", false, err
	}
	if vals["vaild"] != "" {
		hasData = true
	}
	return vals["nick_name"], vals["pic_name"], hasData, nil

}

// SetNickNameAndPicName 设置昵称和头像.
func SetNickNameAndPicName(userName string, nickName string, picName string) error {
	fields := map[string]interface{}{
		"vaild":     "1",
		"nick_name": nickName,
		"pic_name":  picName,
	}
	err := client.HMSet(client.Context(), userName, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

// InvaildCache 将用户数据设置无效，主要用于写入数据库之前，保持数据一直
func InvaildCache(userName string) error {
	err := client.HSet(client.Context(), userName, "vaild", "").Err()
	if err != nil {
		return err
	}
	return nil
}

// SetToken 设置token， 包括token的存活时间
func SetToken(userName string, token string, expiration int64) error {
	err := client.Set(client.Context(), "auth_"+userName, token, time.Duration(expiration*1e9)).Err()
	if err != nil {
		return err
	}
	return nil
}

// CheckToken 校验token
func CheckToken(userName string, token string) (bool, error) {
	val, err := client.Get(client.Context(), "auth_"+userName).Result()
	if err != nil {
		return false, err
	}
	return token == val, nil
}
