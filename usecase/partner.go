package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
)

type PartnerStore struct {
	partnerRepo domain.PartnerRepository
}

func NewPartnerStore(partnerRepo domain.PartnerRepository) *PartnerStore {
	return &PartnerStore{partnerRepo: partnerRepo}
}

func (p PartnerStore) Validation(c *gin.Context) (*domain.Partner, error) {
	var partner domain.Partner
	err := c.ShouldBind(&partner)
	if err != nil {
		return nil, err
	}

	return &partner, err
}

func (p PartnerStore) Store(partner *domain.Partner) error {
	return p.partnerRepo.Store(partner)
}
