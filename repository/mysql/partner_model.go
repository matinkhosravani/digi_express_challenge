package mysql

import (
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"gorm.io/gorm"
)

type partner struct {
	gorm.Model
	TradingName    string
	OwnerName      string
	Document       string `gorm:"unique"`
	Address        address
	CoverageArea   coverageArea
	AddressID      uint
	CoverageAreaID uint
}

type address struct {
	gorm.Model
	Type        string
	Coordinates *location `gorm:"type:POINT"`
}

type coverageArea struct {
	gorm.Model
	Type        string
	Coordinates polygon `gorm:"type:MULTIPOLYGON"`
}

func (p *partner) fromDomain(domainP *domain.Partner) {
	p.OwnerName = domainP.OwnerName
	p.TradingName = domainP.TradingName
	p.Document = domainP.Document
	p.Address = address{
		Type: domainP.Address.Type,
		Coordinates: &location{
			wkb.Point{
				Point: geom.NewPoint(geom.XY).MustSetCoords(domainP.Address.Coordinates).SetSRID(4326),
			},
		},
	}
	p.CoverageArea = coverageArea{
		Type: domainP.CoverageArea.Type,
		Coordinates: polygon{
			MPoly: wkb.MultiPolygon{
				MultiPolygon: geom.NewMultiPolygon(geom.XY).MustSetCoords(convertFloats(domainP.CoverageArea.Coordinates)).SetSRID(4326),
			},
		},
	}
}
func (p *partner) toDomain(domainP *domain.Partner) {

	domainP.ID = p.ID
	domainP.OwnerName = p.OwnerName
	domainP.TradingName = p.TradingName
	domainP.Document = p.Document
	domainP.Address = domain.Address{
		Type:        p.Address.Type,
		Coordinates: []float64{p.Address.Coordinates.Point.X(), p.Address.Coordinates.Point.Y()},
	}
	domainP.CoverageArea = domain.CoverageArea{
		Type:        p.CoverageArea.Type,
		Coordinates: convertCoords(p.CoverageArea.Coordinates.MPoly.MultiPolygon.Coords()),
	}
}
