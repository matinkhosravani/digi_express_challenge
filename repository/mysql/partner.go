package mysql

import (
	"fmt"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PartnerRepository struct {
	DB *gorm.DB
}

func (pr *PartnerRepository) Store(domainP *domain.Partner) error {
	var p partner
	p.fromDomain(domainP)

	err := pr.DB.Transaction(func(tx *gorm.DB) error {
		err := insertAddress(tx, &p)
		if err != nil {
			return err
		}
		p.AddressId = p.address.ID
		err = insertCoverageArea(tx, &p)
		if err != nil {
			return err
		}
		p.CoverageAreaId = p.coverageArea.ID
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
	p.coverageArea.CreatedAt = time.Now()
	p.coverageArea.UpdatedAt = time.Now()
	p.coverageArea.DeletedAt = gorm.DeletedAt{
		Time:  time.Time{},
		Valid: false,
	}
	err := tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`INSERT INTO coverage_areas (created_at,updated_at,deleted_at,type, coordinates)
VALUES (?,?,?,?,  ST_GeomFromText(?));`, p.coverageArea.CreatedAt, p.coverageArea.UpdatedAt, p.coverageArea.DeletedAt, p.coverageArea.Type, generateMultiPolygon(p.coverageArea.Coordinates)).Error
		if err != nil {
			return err
		}
		var insertedID int64
		tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&insertedID)
		if err != nil {
			return err
		}
		p.coverageArea.ID = uint(insertedID)
		return nil
	})

	return err
}

func insertAddress(tx *gorm.DB, p *partner) error {
	p.address.CreatedAt = time.Now()
	p.address.UpdatedAt = time.Now()
	p.address.DeletedAt = gorm.DeletedAt{
		Time:  time.Time{},
		Valid: false,
	}
	err := tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`INSERT INTO addresses (created_at,updated_at,deleted_at,type,coordinates)
VALUES (?,?,NULL,?,POINT(?, ?))`, p.address.CreatedAt, p.address.UpdatedAt, p.address.Type, p.address.Coordinates[0], p.address.Coordinates[1]).Error
		if err != nil {
			return err
		}
		var insertedID int64
		tx.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&insertedID)
		if err != nil {
			return err
		}
		p.address.ID = uint(insertedID)

		return nil
	})

	return err
}

func generateMultiPolygon(coordinates [][][][]float64) string {
	result := "MULTIPOLYGON("
	for _, polygon := range coordinates {
		result += "("
		for _, ring := range polygon {
			result += "("
			for _, point := range ring {
				result += fmt.Sprintf("%.6f %.6f,", point[0], point[1])
			}
			result = strings.TrimSuffix(result, ",") + "),"
		}
		result = strings.TrimSuffix(result, ",") + "),"
	}
	result = strings.TrimSuffix(result, ",") + ")"

	return result
}
