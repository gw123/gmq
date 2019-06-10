package core

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DbPool struct {
	pool    map[string]*gorm.DB
	defualt string
}

func NewDbPool() *DbPool {
	this := new(DbPool)
	this.pool = make(map[string]*gorm.DB, 0)
	return this
}

func (this *DbPool) NewDb(drive ,db_host, db_database, db_username, db_pwd string) (*gorm.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_pwd, db_host, db_database)
	dbInstance, err := gorm.Open(drive, connStr)
	if err != nil {
		return nil, err
	}
	this.pool[db_database] = dbInstance
	return dbInstance, nil
}

func (this *DbPool) SetDb(dbNmae string, db *gorm.DB) {
	if db == nil {
		return
	}
	this.pool[dbNmae] = db
}

func (this *DbPool) GetDb(dbname string) (*gorm.DB, error) {
	db, ok := this.pool[dbname]
	if !ok {
		return nil, errors.New("cant find this db " + dbname)
	}
	return db, nil
}
