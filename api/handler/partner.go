package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/domain"
	"net/http"
	"strconv"
)

type Partner struct {
	StoreUsecase domain.PartnerStoreUsecase
	LoadUsecase  domain.PartnerLoadUsecase
}

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
