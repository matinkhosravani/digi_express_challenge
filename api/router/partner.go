package router

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/digi_express_challenge/api/handler"
	"github.com/matinkhosravani/digi_express_challenge/repository"
	"github.com/matinkhosravani/digi_express_challenge/usecase"
)

func PartnerRoutes(r *gin.RouterGroup) {
	partners := r.Group("/partners")
	pr := repository.NewPartnerRepository()
	h := handler.Partner{
		StoreUsecase: usecase.NewPartnerStore(pr),
		LoadUsecase:  usecase.NewPartnerLoad(pr),
	}

	partners.POST("", h.Store)
	partners.GET(":id", h.LoadByID)
}
