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
	"net/url"
	"strconv"
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

func TestUnitPartner_LoadPartner(t *testing.T) {
	p := domain.Partner{
		ID:          1,
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

	//happy path
	t.Run("User loads a partner with id", func(t *testing.T) {
		u := mock.PartnerLoadUsecase{
			GetPartnerByIdFn: func(ID uint) (*domain.Partner, error) {
				return &p, nil
			},
		}

		w := loadPartner(&u, &p)
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

	t.Run("User loads a partner with  invalid id", func(t *testing.T) {
		dummyError := fmt.Errorf("dummy")
		u := mock.PartnerLoadUsecase{
			GetPartnerByIdFn: func(ID uint) (*domain.Partner, error) {
				return &p, dummyError
			},
		}

		w := loadPartner(&u, &p)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
		assert.Equal(t, dummyError.Error(), response["error"])

	})

}

func TestUnitPartner_SearchPartner(t *testing.T) {
	p := domain.Partner{
		ID:          1,
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

	//happy path
	t.Run("User search partner with Point", func(t *testing.T) {
		u := mock.PartnerSearchUsecase{
			SearchPartnersFn: func(x, y float64, limit int) ([]*domain.Partner, error) {
				p.ID = 1
				return []*domain.Partner{&p}, nil
			},
			ValidationFn: func(c *gin.Context) (*domain.PartnerSearchRequest, error) {
				return &domain.PartnerSearchRequest{
					X: 20,
					Y: 20,
				}, nil
			},
		}

		w := searchPartner(&u, domain.PartnerSearchRequest{
			X: 20,
			Y: 20,
		})
		assert.Equal(t, http.StatusOK, w.Code)
		var responsePartner []domain.Partner
		_ = json.Unmarshal(w.Body.Bytes(), &responsePartner)
		assert.Equal(t, responsePartner[0].ID, uint(1))
		assert.Equal(t, responsePartner[0].TradingName, p.TradingName)
		assert.Equal(t, responsePartner[0].OwnerName, p.OwnerName)
		assert.Equal(t, responsePartner[0].Document, p.Document)
		assert.Equal(t, responsePartner[0].Address, p.Address)
		assert.Equal(t, responsePartner[0].CoverageArea, p.CoverageArea)
	})

	t.Run("User gets error with invalid request", func(t *testing.T) {
		dummyError := fmt.Errorf("dummy")
		u := mock.PartnerSearchUsecase{
			ValidationFn: func(c *gin.Context) (*domain.PartnerSearchRequest, error) {
				return nil, dummyError
			},
		}

		w := searchPartner(&u, domain.PartnerSearchRequest{})
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotZero(t, response["error"])
		assert.Equal(t, dummyError.Error(), response["error"])

	})

}
func loadPartner(u domain.PartnerLoadUsecase, p *domain.Partner) *httptest.ResponseRecorder {
	h := &Partner{
		LoadUsecase: u,
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: strconv.Itoa(int(p.ID)),
		},
	}
	h.LoadByID(c)

	return w
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

func searchPartner(u domain.PartnerSearchUsecase, request domain.PartnerSearchRequest) *httptest.ResponseRecorder {
	h := &Partner{
		SearchUsecase: u,
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	v := url.Values{}
	v.Add("x", fmt.Sprintf("%v", request.X))
	v.Add("y", fmt.Sprintf("%v", request.Y))
	c.Request = &http.Request{
		Header: make(http.Header),
		Form:   v,
	}
	h.Search(c)

	return w
}
