package services

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/redisKeys"
	"github.com/gw123/GMQ/common/utils"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type RegisterParam struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type LoginParam struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

type CollectionItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type RegisterUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Password string `json:"password2"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (r *RegisterUser) TableName() string {
	return "users"
}

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

func (s *UserService) GetUser(id int) (*User, error) {
	var user User
	_, err := utils.GetCache(s.redis, redisKeys.User+strconv.Itoa(id), &user, func() (interface{}, error) {
		res := s.db.Where("id = ?", id).Find(&user)
		if res.Error != nil {
			return nil, res.Error
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Register(params RegisterParam) error {
	//判断用户是否存在
	res := s.redis.Get(redisKeys.MessageCheckCode + params.Mobile)
	if res.Val() != params.Code {
		return errors.New("验证码校验失败")
	}

	user := &RegisterUser{}
	if res := s.db.Where("mobile = ?", params.Mobile).Find(user); res.Error != nil && !res.RecordNotFound() {
		return res.Error
	} else if !res.RecordNotFound() {
		return errors.New("手机号码已经注册")
	}

	//
	key, err := s.app.GetAppConfigItem("appKey")
	if err != nil {
		return err
	}

	newPwd, err := utils.MakePassword(key, params.Password)
	if err != nil {
		return err
	}
	newUser := &RegisterUser{
		Name:     params.Mobile,
		Mobile:   params.Mobile,
		Password: newPwd,
		Email:    "",
	}

	if err := s.db.Save(newUser).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(params LoginParam) (string, error) {
	//s.db.LogMode(true)
	key, err := s.app.GetAppConfigItem("appKey")
	if err != nil {
		return "", err
	}
	newPwd, err := utils.MakePassword(key, params.Password)
	if err != nil {
		return "", err
	}
	user := &RegisterUser{}
	if err := s.db.Where("mobile = ?", params.Mobile).Find(user).Error; err != nil {
		return "", err
	}
	if user.Password != newPwd {
		return "", errors.New("账号或者密码错误")
	}

	info := map[string]string{
		"name":   user.Name,
		"mobile": user.Mobile,
		"avatar": user.Avatar,
	}
	res, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	claims := jwt.StandardClaims{
		Audience:  "mobile",
		ExpiresAt: time.Now().Unix() + 3600*24*7,
		Id:        strconv.Itoa(user.Id),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "xytschool",
		Subject:   string(res),
	}

	jwtKey, ok := s.app.GetConfig().Get("modules.web.jwt_key").(string)
	if !ok {
		return "", errors.New("配置问题")
	}

	token, err := utils.MakeJwtTokenSh1(&claims, jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

type UserCollection struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	TargetId  int       `json:"target_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UserCollection) TableName() string {
	return "user_collection"
}

func (s *UserService) AddUserCollection(userId int, targetType string, targetId int) (error) {
	//s.db.LogMode(true)
	collecion := &UserCollection{}

	if !s.db.Where("user_id = ?", userId).
		Where("type = ?", targetType).
		Where("target_id = ?", targetId).Find(collecion).RecordNotFound() {
		return errors.New("已经收藏")
	}

	collecion.UserId = userId
	collecion.Type = targetType
	collecion.TargetId = targetId
	collecion.CreatedAt = time.Now()
	if err := s.db.Save(collecion).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) DelUserCollection(userId int, targetType string, targetId int) (error) {
	collecion := &UserCollection{}
	collecion.UserId = userId
	collecion.Type = targetType
	collecion.TargetId = targetId

	if err := s.db.Where("user_id = ?", userId).
		Where("type = ?", targetType).
		Where("target_id = ?", targetId).
		Delete(collecion).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserCollection(userId int) ([]CollectionItem, error) {
	var items []CollectionItem
	//s.db.LogMode(true)
	res := s.db.Select("g.id,g.title").Table("user_collection as  c").
		Joins("left join `group` as g on g.id =c.target_id ").
		Where("c.user_id = ?", userId).
		Limit(20).
		Order("c.id desc").
		Find(&items)

	if res.Error != nil && !res.RecordNotFound() {
		return nil, res.Error
	}
	return items, nil
}
