package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/app"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"github.com/matinkhosravani/digi_express_challenge/domain/factory"
	"github.com/matinkhosravani/digi_express_challenge/repository"
	"github.com/matinkhosravani/digi_express_challenge/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestIntegrationPartner_StorePartner(t *testing.T) {
	app.BootTestApp()
	pr := repository.NewPartnerRepository()
	uc := usecase.NewPartnerStore(pr)

	//happy path
	t.Run("User can store valid partner into database", func(t *testing.T) {
		pr.Empty()
		p := factory.NewDefaultPartnerFactory(pr).Get()
		w := storePartner(uc, p)
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, responsePartner.ID, p.ID)
		assert.Equal(t, responsePartner.TradingName, p.TradingName)
		assert.Equal(t, responsePartner.OwnerName, p.OwnerName)
		assert.Equal(t, responsePartner.Document, p.Document)
		assert.Equal(t, responsePartner.Address, p.Address)
		assert.Equal(t, responsePartner.CoverageArea, p.CoverageArea)
	})

	t.Run("User cant store a partner with duplicate document", func(t *testing.T) {
		pr.Empty()
		//store a partner in database
		p := factory.NewDefaultPartnerFactory(pr).Build()
		p2 := factory.NewDefaultPartnerFactory(pr).
			WithDocument(p.Document). //setting new partner's document to existing partner's document
			Get()
		w := storePartner(uc, p2)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
	})
}

func TestIntegrationPartner_LoadPartner(t *testing.T) {
	app.BootTestApp()
	pr := repository.NewPartnerRepository()
	uc := usecase.NewPartnerLoad(pr)

	//happy path
	t.Run("User loads a partner with id", func(t *testing.T) {
		pr.Empty()
		p := factory.NewDefaultPartnerFactory(pr).Build()
		w := loadPartner(uc, p)
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, responsePartner.ID, p.ID)
		assert.Equal(t, responsePartner.TradingName, p.TradingName)
		assert.Equal(t, responsePartner.OwnerName, p.OwnerName)
		assert.Equal(t, responsePartner.Document, p.Document)
		assert.Equal(t, responsePartner.Address, p.Address)
		assert.Equal(t, responsePartner.CoverageArea, p.CoverageArea)
	})
	//happy path
	t.Run("User gets error when loading a partner that does not exist in database", func(t *testing.T) {
		pr.Empty()
		h := &Partner{
			LoadUsecase: uc,
		}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: strconv.Itoa(int(777)), //dummy id
			},
		}
		h.LoadByID(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
	})

	t.Run("User loads a partner with invalid id", func(t *testing.T) {
		g := gin.Default()
		rec := httptest.NewRecorder()
		h := Partner{LoadUsecase: uc}
		g.GET("/partners/:id", h.LoadByID)
		req := httptest.NewRequest(http.MethodGet, "/partners/invalid", nil)
		g.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		_ = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
	})

}

func TestIntegrationPartner_SearchPartner(t *testing.T) {
	app.BootTestApp()
	pr := repository.NewPartnerRepository()
	uc := usecase.NewPartnerSearch(pr)
	//happy path
	t.Run("User search partner with Point", func(t *testing.T) {
		pr.Empty()
		partner1 := factory.NewDefaultPartnerFactory(pr).
			WithCoverageArea(domain.CoverageArea{
				Type: "MULTIPOLYGON",
				Coordinates: [][][][]float64{
					{
						{ //sqaure
							{0, 0},
							{40, 0},
							{40, 40},
							{0, 40},
							{0, 0},
						},
					},
				},
			}).
			WithAddress(domain.Address{
				Type:        "POINT",
				Coordinates: []float64{20, 0},
			}).Build()
		_ = factory.NewDefaultPartnerFactory(pr).
			WithCoverageArea(domain.CoverageArea{
				Type: "MULTIPOLYGON",
				Coordinates: [][][][]float64{
					{
						{ //diamond
							{20, 0},
							{40, 20},
							{20, 40},
							{0, 20},
							{20, 0},
						},
					},
				},
			}).
			WithAddress(domain.Address{
				Type:        "POINT",
				Coordinates: []float64{0, 0},
			}).Build()
		//parcel is at center of diamond
		w := searchPartner(uc, domain.PartnerSearchRequest{
			X: 20,
			Y: 20,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner []domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, responsePartner[0].ID, partner1.ID)
		assert.Equal(t, responsePartner[0].TradingName, partner1.TradingName)
		assert.Equal(t, responsePartner[0].OwnerName, partner1.OwnerName)
		assert.Equal(t, responsePartner[0].Document, partner1.Document)
		assert.Equal(t, responsePartner[0].Address, partner1.Address)
		assert.Equal(t, responsePartner[0].CoverageArea, partner1.CoverageArea)
	})

	t.Run("Search partner does not get a partner which is close to parcel but parcel is not in his coverage area", func(t *testing.T) {
		pr.Empty()
		_ = factory.NewDefaultPartnerFactory(pr).
			WithCoverageArea(domain.CoverageArea{
				Type: "MULTIPOLYGON",
				Coordinates: [][][][]float64{
					{
						{ //sqaure
							{0, 0},
							{40, 0},
							{40, 40},
							{0, 40},
							{0, 0},
						},
					},
				},
			}).
			WithAddress(domain.Address{
				Type:        "POINT",
				Coordinates: []float64{50, 50},
			}).Build()

		//parcel is at center of diamond
		w := searchPartner(uc, domain.PartnerSearchRequest{
			X: 50,
			Y: 50,
		})
		//partner1 is closer
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner []domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, 0, len(responsePartner))
	})
	t.Run("Search partner does not get a partner when parcel is on the edge of coverage area", func(t *testing.T) {
		pr.Empty()
		_ = factory.NewDefaultPartnerFactory(pr).
			WithCoverageArea(domain.CoverageArea{
				Type: "MULTIPOLYGON",
				Coordinates: [][][][]float64{
					{
						{ //sqaure
							{0, 0},
							{40, 0},
							{40, 40},
							{0, 40},
							{0, 0},
						},
					},
				},
			}).
			WithAddress(domain.Address{
				Type:        "POINT",
				Coordinates: []float64{50, 50},
			}).Build()

		w := searchPartner(uc, domain.PartnerSearchRequest{
			X: 40,
			Y: 40,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner []domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, 0, len(responsePartner))

		w = searchPartner(uc, domain.PartnerSearchRequest{
			X: 39,
			Y: 40,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, 0, len(responsePartner))
	})
}
