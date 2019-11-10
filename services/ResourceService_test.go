package services

import (
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/models"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"testing"
)

func TestResourceService_SaveGroup(t *testing.T) {
	type fields struct {
		app   interfaces.App
		db    *gorm.DB
		redis *redis.Client
	}
	type args struct {
		ctx   echo.Context
		group *models.Group
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ResourceService{
				app:   tt.fields.app,
				db:    tt.fields.db,
				redis: tt.fields.redis,
			}
			if err := s.SaveGroup(tt.args.ctx, tt.args.group); (err != nil) != tt.wantErr {
				t.Errorf("SaveGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}