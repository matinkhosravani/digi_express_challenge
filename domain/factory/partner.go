package factory

import (
	"github.com/bxcodec/faker/v4"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"log"
)

type PartnerFactory struct {
	partnerRepo domain.PartnerRepository
	partner     *domain.Partner
}

// NewDefaultPartnerFactory Returning Pointer to struct for chaining methods
// it's kinda a form of Builder pattern
func NewDefaultPartnerFactory(partnerRepo domain.PartnerRepository) *PartnerFactory {
	pf := &PartnerFactory{
		partnerRepo: partnerRepo,
		partner: &domain.Partner{
			ID:          1,
			OwnerName:   faker.Name(),
			Document:    faker.UUIDDigit(),
			TradingName: faker.Name(),
			Address: domain.Address{
				Type:        "POINT",
				Coordinates: []float64{0, 0},
			},
			CoverageArea: domain.CoverageArea{
				Type: "MULTIPOLYGON",
				Coordinates: [][][][]float64{
					{
						{
							{0, 0},
							{40, 0},
							{40, 40},
							{0, 40},
							{0, 0},
						},
					},
				},
			},
		},
	}

	return pf
}
func (pf *PartnerFactory) WithID(id uint) *PartnerFactory {
	pf.partner.ID = id
	return pf
}

func (pf *PartnerFactory) WithOwnerName(n string) *PartnerFactory {
	pf.partner.OwnerName = n
	return pf
}
func (pf *PartnerFactory) WithDocument(n string) *PartnerFactory {
	pf.partner.Document = n
	return pf
}
func (pf *PartnerFactory) WithTradingName(n string) *PartnerFactory {
	pf.partner.TradingName = n
	return pf
}

func (pf *PartnerFactory) WithAddress(a domain.Address) *PartnerFactory {
	pf.partner.Address = a
	return pf
}

func (pf *PartnerFactory) WithCoverageArea(ca domain.CoverageArea) *PartnerFactory {
	pf.partner.CoverageArea = ca
	return pf
}
func (pf *PartnerFactory) Get() *domain.Partner {
	return pf.partner
}
func (pf *PartnerFactory) Build() *domain.Partner {
	err := pf.partnerRepo.Store(pf.partner)
	if err != nil {
		log.Fatal(err)
	}

	return pf.partner
}
