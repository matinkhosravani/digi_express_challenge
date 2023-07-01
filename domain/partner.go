package domain

import (
	"github.com/gin-gonic/gin"
)

type Partner struct {
	ID           uint         `json:"id" swaggerignore:"true"`
	TradingName  string       `json:"tradingName" example:"Adega da Cerveja - Pinheiros"`
	OwnerName    string       `json:"ownerName" example:"Zé da Silva"`
	Document     string       `json:"document" example:"1432132123891/0001"`
	CoverageArea CoverageArea `json:"coverageArea"`
	Address      Address      `json:"address"`
}
type Address struct {
	Type        string    `json:"type" example:"Point"`
	Coordinates []float64 `json:"coordinates"`
}

type CoverageArea struct {
	Type        string          `json:"type" example:"MultiPolygon"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type Repository interface {
	Empty() error
}

type PartnerRepository interface {
	Repository
	Store(*Partner) error
	GetByID(ID uint) (*Partner, error)
	SearchPartners(x, y float64, limit int) ([]*Partner, error)
}

type PartnerStoreUsecase interface {
	Validation(c *gin.Context) (*Partner, error)
	Store(partner *Partner) error
}

type PartnerLoadUsecase interface {
	GetPartnerById(ID uint) (*Partner, error)
}

type PartnerSearchUsecase interface {
	Validation(c *gin.Context) (*PartnerSearchRequest, error)
	SearchPartners(x, y float64, limit int) ([]*Partner, error)
}

type PartnerSearchRequest struct {
	X float64 `form:"x" binding:"required,latitude"`
	Y float64 `form:"y" binding:"required,longitude"`
}
