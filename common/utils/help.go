package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"regexp"
	"time"
)

var reg = regexp.MustCompile(`^\$\{(.*)\}$`)

//判断配置是否是环境变量 如${evn_name} 类型
func IsEvnParam(key string) bool {
	if reg.MatchString(key) {
		arrs := reg.FindStringSubmatch(key)
		if len(arrs) > 1 {
			return true
		}
	}
	return false
}

//从环境变量加载
func LoadFromEnv(key string) string {
	if reg.MatchString(key) {
		arrs := reg.FindStringSubmatch(key)
		if len(arrs) > 1 {
			return os.Getenv(arrs[1])
		}
	}
	return ""
}

//先从缓存读取数据,如果不存在调用 call方法获取后在存放到数据库中
func GetCache(client *redis.Client, key string, out interface{}, call func() (interface{}, error)) (string, error) {
	val := client.Get(key).Val()
	//是否需要重新刷新缓存
	isOK := false
	if val != "" {
		if out != nil {
			err := json.Unmarshal([]byte(val), out)
			if err != nil {
				isOK = false
			} else {
				isOK = true
			}
		} else {
			isOK = true
		}
	}

	if val == "" || !isOK {
		newVal, err := call()
		if err != nil {
			return "", err
		}
		tmp, err := json.Marshal(newVal)
		if err != nil {
			return "", err
		}
		if err := client.Set(key, tmp, time.Hour).Err(); err != nil {
			return "", err
		}
		if err := json.Unmarshal([]byte(tmp), out); err != nil {
			return "", err
		}
		return string(tmp), err
	}

	return val, nil
}

//make passwd

func MakePassword(key, rawpwd string) (string, error) {
	pwdEncoder := hmac.New(sha1.New, []byte(key))
	_, err := pwdEncoder.Write([]byte(rawpwd))
	if err != nil {
		return "", err
	}
	passwd := pwdEncoder.Sum(nil)
	return fmt.Sprintf("%x", passwd), nil
}
