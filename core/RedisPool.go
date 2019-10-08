package core

import (
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/utils"
	"github.com/gw123/GMQ/core/interfaces"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type RedisPool struct {
	pool    map[string]*redis.Client
	defualt string
	config  map[string]interface{}
	app     interfaces.App
}

func NewRedisPool(config map[string]interface{}, app interfaces.App) *RedisPool {
	this := new(RedisPool)
	this.pool = make(map[string]*redis.Client)
	this.config = config
	this.app = app
	this.loadConfig()
	return this
}

func (this *RedisPool) loadConfig() {
	configs := this.config
	for key, config := range configs {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			if key != "default" {
				this.app.Info("App", "redis配置文件 格式不支持")
			}
			continue
		}

		host, ok := configMap["host"].(string)
		if !ok {
			host = "127.0.0.1"
		} else if utils.IsEvnParam(host) {
			host = utils.LoadFromEnv(host)
		}

		port, ok := configMap["port"].(string)
		if !ok || port == "" {
			port = "6379"
		} else if utils.IsEvnParam(port) {
			port = utils.LoadFromEnv(port)
		}

		database, ok := configMap["database"].(int)
		if !ok {
			database = 0
		}

		password, ok := configMap["password"].(string)
		if !ok {
			password = ""
		} else if utils.IsEvnParam(password) {
			password = utils.LoadFromEnv(password)
		}
		this.app.Info("App", "load redis name:%s ,database:%s", key, database)
		db, err := this.NewRedis(
			host,
			port,
			password,
			database)
		if err != nil {
			this.app.Warning("App", "redis  load error, %d: ", err.Error())
		}
		this.SetDb(key, db)
	}

	//set default db
	defaultDBkey, ok := configs["default"].(string)
	if !ok {
		return
	}
	this.app.Debug("App", "default Redis: %s", defaultDBkey)

	db, err := this.GetDb(defaultDBkey)
	if err != nil {
		this.app.Warning("App", "cant found redis default config :%s", err.Error())
	} else {
		this.SetDb("default", db)
	}
}

func (this *RedisPool) NewRedis(host, port, pwd string, database int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pwd,
		DB:       database,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (this *RedisPool) SetDb(dbNmae string, db *redis.Client) {
	if db == nil {
		return
	}
	this.pool[dbNmae] = db
}

func (this *RedisPool) GetDb(dbname string) (*redis.Client, error) {
	db, ok := this.pool[dbname]
	if !ok {
		return nil, errors.New("cant find this db " + dbname)
	}
	return db, nil
}
