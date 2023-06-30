package domain

import "github.com/gin-gonic/gin"

type Partner struct {
	ID           uint         `json:"id"`
	TradingName  string       `json:"tradingName"`
	OwnerName    string       `json:"ownerName"`
	Document     string       `json:"document"`
	CoverageArea CoverageArea `json:"coverageArea"`
	Address      Address      `json:"address"`
}
type Address struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type CoverageArea struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type PartnerRepository interface {
	Store(*Partner) error
	GetByID(ID uint) (*Partner, error)
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
