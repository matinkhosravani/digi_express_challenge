package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
)

type PartnerStoreUsecase struct {
	ValidationFn func(c *gin.Context) (*domain.Partner, error)
	StoreFn      func(partner *domain.Partner) error
}

func (uc *PartnerStoreUsecase) Validation(c *gin.Context) (*domain.Partner, error) {
	if uc != nil && uc.ValidationFn != nil {
		return uc.ValidationFn(c)
	}

	var p domain.Partner
	err := c.BindJSON(&p)
	return &p, err
}

func (uc *PartnerStoreUsecase) Store(partner *domain.Partner) error {
	if uc != nil && uc.StoreFn != nil {
		return uc.StoreFn(partner)
	}

	return nil
}

type PartnerLoadUsecase struct {
	GetPartnerByIdFn func(ID uint) (*domain.Partner, error)
}

func (p *PartnerLoadUsecase) GetPartnerById(ID uint) (*domain.Partner, error) {
	if p != nil && p.GetPartnerByIdFn != nil {
		return p.GetPartnerByIdFn(ID)
	}

	return &domain.Partner{}, nil
}
