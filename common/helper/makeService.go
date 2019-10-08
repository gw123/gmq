package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var serviceContent = `package services

import (
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/redisKeys"
	"github.com/gw123/GMQ/common/utils"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
)


type UserService struct {
	app   interfaces.App
	db    *gorm.DB
	redis *redis.Client
}

func NewUserService(app interfaces.App) (*UserService, error) {
	db, err := app.GetDefaultDb()
	if err != nil {
		return nil, err
	}

	redisClient, err := app.GetDefaultRedis()
	if err != nil {
		return nil, err
	}

	return &UserService{
		app:   app,
		db:    db,
		redis: redisClient,
	}, nil
}

func (s *UserService) GetServiceName() string {
	return "UserService"
}
`

func MakeService(moduleName string, distDir string) error {
	if distDir == "" {
		distDir = "./services/"
	}

	if !strings.HasSuffix(distDir, "/") {
		distDir += "/"
	}
	filePath := distDir + moduleName

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && fileInfo.IsDir() {
		return errors.New("dir is exist")
	}

	fmt.Println("Mkdir :" + filePath)
	err = os.MkdirAll(filePath, 0660)
	if err != nil {
		return err
	}

	f, err := os.Create(filePath + "/" + moduleName + ".go")
	if err != nil {
		return err
	}
	defer f.Close()
	content := MakeServiceContent(moduleName)
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func MakeServiceContent(serviceName string) string {
	return strings.Replace(serviceContent, "UserService", serviceName, -1)
}
