package repository

import (
	"github.com/matinkhosravani/digi_express_challenge/app"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"github.com/matinkhosravani/digi_express_challenge/repository/mysql"
	"log"
)

func NewPartnerRepository() domain.PartnerRepository {
	switch app.GetEnv().DBType {
	case "mysql":
		db, err := mysql.NewGorm()
		if err != nil {
			log.Fatal(err.Error())
		}
		return &mysql.PartnerRepository{DB: db}
	}

	return nil
}
