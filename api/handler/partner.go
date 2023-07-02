package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"net/http"
	"strconv"
)

type Partner struct {
	StoreUsecase  domain.PartnerStoreUsecase
	LoadUsecase   domain.PartnerLoadUsecase
	SearchUsecase domain.PartnerSearchUsecase
}

// Store
// @Summary Store Partner
// @Description Store a partner
// @ID partner-store
// @Accept json
// @Produce json
// @Param req body domain.Partner true "Partner object"
// @Success 200 {object} domain.Partner
// @Router /partners/ [post]
func (h *Partner) Store(c *gin.Context) {
	partner, err := h.StoreUsecase.Validation(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = h.StoreUsecase.Store(partner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, partner)
}

// LoadByID is the handler function for checking if a username is taken.
// @Summary Load a partner
// @Description Load a partner with the id in the url
// @ID partner-load
// @Produce json
// @Param id path int true "id of partner"
// @Success 200 {object} domain.Partner
// @Router /partners/{id} [get]
func (h *Partner) LoadByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	partner, err := h.LoadUsecase.GetPartnerById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, partner)
}

// Search is the handler function for finding the nearest partner to a Point.
// @Summary Find the nearest Partner
// @Description is the handler function for finding the nearest partner to a Point.
// @ID partner-search
// @Produce json
// @Param x query number true "x of point"
// @Param y query number true "y of point"
// @Success 200 {object} domain.Partner
// @Router /partners/search [get]
func (h *Partner) Search(c *gin.Context) {
	request, err := h.SearchUsecase.Validation(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	partner, err := h.SearchUsecase.SearchPartners(request.X, request.Y, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, partner)
}
