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

func (pr *PartnerRepository) Empty() error {
	err := pr.DB.Exec("drop Table partners").Error
	if err != nil {
		return err
	}
	err = pr.DB.Exec("TRUNCATE Table coverage_areas").Error
	if err != nil {
		return err
	}
	err = pr.DB.Exec("Truncate Table addresses").Error
	if err != nil {
		return err
	}

	err = pr.DB.AutoMigrate(&partner{})
	if err != nil {
		return err
	}
	return err
}

func (pr *PartnerRepository) SearchPartners(x, y float64, limit int) ([]*domain.Partner, error) {
	var ps []partner
	point := fmt.Sprintf("POINT(%v,%v)", x, y)
	err := pr.DB.Table("partners").
		Select("partners.*").
		Joins("JOIN coverage_areas ca ON partners.coverage_area_id = ca.id").
		Joins("JOIN addresses a ON a.id = partners.address_id").Preload("Address").
		Where(fmt.Sprintf("ST_Within(%s, ca.coordinates)", point)).
		Order(fmt.Sprintf("ST_Distance(%s, a.coordinates)", point)).
		Limit(limit).
		Preload("CoverageArea").
		Preload("Address").
		Find(&ps).Error

	if err != nil {
		return nil, err
	}

	domainPs := make([]*domain.Partner, 0, len(ps))
	for _, p := range ps {
		var domainP domain.Partner
		p.toDomain(&domainP)
		domainPs = append(domainPs, &domainP)
	}

	return domainPs, nil
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
