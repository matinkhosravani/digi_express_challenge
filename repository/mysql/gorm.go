package mysql

import (
	"fmt"
	"github.com/matinkhosravani/digi_express_challenge/app"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGorm() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		app.GetEnv().DBUser,
		app.GetEnv().DBPass,
		app.GetEnv().DBHost,
		app.GetEnv().DBPort,
		app.GetEnv().DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&address{}, &coverageArea{}, &partner{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
