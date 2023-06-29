package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
)

type PartnerStore struct {
	partnerRepo domain.PartnerRepository
}

func NewPartnerStore(partnerRepo domain.PartnerRepository) domain.PartnerStoreUsecase {
	return &PartnerStore{partnerRepo: partnerRepo}
}
func NewPartnerLoad(partnerRepo domain.PartnerRepository) domain.PartnerLoadUsecase {
	return &PartnerLoad{partnerRepo: partnerRepo}
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

type PartnerLoad struct {
	partnerRepo domain.PartnerRepository
}

func (p PartnerLoad) GetPartnerById(ID uint) (*domain.Partner, error) {
	return p.partnerRepo.GetByID(ID)
}
