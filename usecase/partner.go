package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
)

type PartnerStore struct {
	partnerRepo domain.PartnerRepository
}

type PartnerLoad struct {
	partnerRepo domain.PartnerRepository
}

type PartnerSearch struct {
	partnerRepo domain.PartnerRepository
}

func NewPartnerSearch(partnerRepo domain.PartnerRepository) *PartnerSearch {
	return &PartnerSearch{partnerRepo: partnerRepo}
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

func (p PartnerLoad) GetPartnerById(ID uint) (*domain.Partner, error) {
	return p.partnerRepo.GetByID(ID)
}

func (p PartnerSearch) Validation(c *gin.Context) (*domain.PartnerSearchRequest, error) {
	var request domain.PartnerSearchRequest
	err := c.ShouldBind(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (p PartnerSearch) SearchPartners(x, y float64, limit int) ([]*domain.Partner, error) {
	return p.partnerRepo.SearchPartners(x, y, limit)
}
