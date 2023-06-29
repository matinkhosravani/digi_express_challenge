package mysql

import (
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"gorm.io/gorm"
)

type partner struct {
	gorm.Model
	TradingName    string
	OwnerName      string
	Document       string `gorm:"unique"`
	address        address
	coverageArea   coverageArea
	AddressId      uint
	CoverageAreaId uint
}

type address struct {
	gorm.Model
	Type        string
	Coordinates []float64 `gorm:"type:POINT"`
	Partners    []partner
}
type coverageArea struct {
	gorm.Model
	Type        string
	Coordinates [][][][]float64 `gorm:"type:MULTIPOLYGON"`
	Partners    []partner
}

func (p *partner) fromDomain(domainP *domain.Partner) {
	p.OwnerName = domainP.OwnerName
	p.TradingName = domainP.TradingName
	p.Document = domainP.Document
	p.address = address{
		Type:        domainP.Address.Type,
		Coordinates: domainP.Address.Coordinates,
	}
	p.coverageArea = coverageArea{
		Type:        domainP.CoverageArea.Type,
		Coordinates: domainP.CoverageArea.Coordinates,
	}
}
func (p *partner) toDomain(domainP *domain.Partner) {
	domainP.ID = p.ID
	domainP.OwnerName = p.OwnerName
	domainP.TradingName = p.TradingName
	domainP.Document = p.Document
	domainP.Address = domain.Address{
		Type:        p.address.Type,
		Coordinates: p.address.Coordinates,
	}
	domainP.CoverageArea = domain.CoverageArea{
		Type:        p.coverageArea.Type,
		Coordinates: p.coverageArea.Coordinates,
	}
}
