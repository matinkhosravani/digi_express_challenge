package mysql

import (
	"fmt"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"gorm.io/gorm"
	"time"
)

type PartnerRepository struct {
	DB *gorm.DB
}

func (pr *PartnerRepository) GetByID(ID uint) (*domain.Partner, error) {
	var p partner
	res := pr.DB.Where("partners.id = ?", ID).
		Joins("JOIN addresses ON addresses.id = partners.address_id").Preload("Address").
		Joins("JOIN coverage_areas ON coverage_areas.id = partners.coverage_area_id").Preload("CoverageArea").
		First(&p)

	if err := res.Error; err != nil {
		return nil, err
	}
	var pDomain domain.Partner
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no such user; id : %d", ID)
	}
	p.toDomain(&pDomain)
	return &pDomain, nil
}

func (pr *PartnerRepository) Store(domainP *domain.Partner) error {
	var p partner
	p.fromDomain(domainP)

	err := pr.DB.Transaction(func(tx *gorm.DB) error {
		err := insertAddress(tx, &p)
		if err != nil {
			return err
		}
		p.AddressID = p.Address.ID
		err = insertCoverageArea(tx, &p)
		if err != nil {
			return err
		}
		p.CoverageAreaID = p.CoverageArea.ID
		err = tx.Create(&p).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	p.toDomain(domainP)

	return nil
}

func insertCoverageArea(tx *gorm.DB, p *partner) error {
	p.CoverageArea.CreatedAt = time.Now()
	p.CoverageArea.UpdatedAt = time.Now()
	p.CoverageArea.DeletedAt = gorm.DeletedAt{
		Time:  time.Time{},
		Valid: false,
	}
	err := tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&p.CoverageArea).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func insertAddress(tx *gorm.DB, p *partner) error {
	p.Address.CreatedAt = time.Now()
	p.Address.UpdatedAt = time.Now()
	p.Address.DeletedAt = gorm.DeletedAt{
		Time:  time.Time{},
		Valid: false,
	}
	err := tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&p.Address).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

