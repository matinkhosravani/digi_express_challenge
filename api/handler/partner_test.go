package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"github.com/matinkhosravani/digi_express_challenge/usecase/mock"
	"github.com/matinkhosravani/digi_express_challenge/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnitPartner_StorePartner(t *testing.T) {

	//happy path
	t.Run("User can store valid partner", func(t *testing.T) {
		p := domain.Partner{
			TradingName: "Adega da Cerveja - Pinheiros",
			OwnerName:   "Zé da Silva",
			Document:    "1432132123891/0003",
			CoverageArea: domain.CoverageArea{
				Type: "MultiPolygon",
				Coordinates: [][][][]float64{
					{
						{
							{30, 20},
							{45, 40},
							{10, 40},
							{30, 20},
						},
					},
					{
						{
							{15, 5},
							{40, 10},
							{10, 20},
							{5, 10},
							{15, 5},
						},
					},
				},
			},
			Address: domain.Address{
				Type:        "Point",
				Coordinates: []float64{-46.57421, -21.785741},
			},
		}
		u := mock.PartnerStoreUsecase{
			StoreFn: func(partner *domain.Partner) error {
				partner.ID = 1 //setting a fake id
				return nil
			},
		}

		w := storePartner(&u, &p)
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, responsePartner.ID, uint(1))
		assert.Equal(t, responsePartner.TradingName, p.TradingName)
		assert.Equal(t, responsePartner.OwnerName, p.OwnerName)
		assert.Equal(t, responsePartner.Document, p.Document)
		assert.Equal(t, responsePartner.Address, p.Address)
		assert.Equal(t, responsePartner.CoverageArea, p.CoverageArea)
	})

	t.Run("User gets error on storing invalid user", func(t *testing.T) {
		p := domain.Partner{
			TradingName: "Adega da Cerveja - Pinheiros",
			OwnerName:   "Zé da Silva",
			Document:    "1432132123891/0003",
			CoverageArea: domain.CoverageArea{
				Type: "MultiPolygon",
				Coordinates: [][][][]float64{
					{
						{
							{30, 20},
							{45, 40},
							{10, 40},
							{30, 20},
						},
					},
					{
						{
							{15, 5},
							{40, 10},
							{10, 20},
							{5, 10},
							{15, 5},
						},
					},
				},
			},
			Address: domain.Address{
				Type:        "",
				Coordinates: []float64{-46.57421, -21.785741},
			},
		}
		dummyError := fmt.Errorf("duumy error")
		u := mock.PartnerStoreUsecase{
			ValidationFn: func(c *gin.Context) (*domain.Partner, error) {
				return &domain.Partner{}, dummyError
			},
		}

		w := storePartner(&u, &p)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
		assert.Equal(t, dummyError.Error(), response["error"])
	})
}

func storePartner(u domain.PartnerStoreUsecase, p *domain.Partner) *httptest.ResponseRecorder {
	h := &Partner{
		StoreUsecase: u,
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	util.MockJsonPost(c, p)

	h.Store(c)

	return w
}
