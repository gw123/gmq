package gmq2

import "github.com/jinzhu/gorm"

type DbPool interface {
	SetDb(dbNmae string, db *gorm.DB)
	GetDb(dbname string) (*gorm.DB, error)
	NewDb(drive, db_host, db_port, db_database, db_username, db_pwd string) (*gorm.DB, error)
}
