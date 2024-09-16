package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"testovoe.com/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		api.GET("/ping", handler.ping)
		api.GET("/tenders", handler.GetTenders)
		api.GET("/bids:tenderId/list", handler.GetBidsByTenderId)
		tenders := api.Group("/tenders")
		{
			tenders.POST("/new", handler.CreateTender)
			tenders.GET("/my", handler.GetTendersByUsername)
			tenders.GET(":tenderId/status", handler.GetTenderStatus)
			tenders.PATCH(":tenderId/edit", handler.EditTender)
			tenders.PUT(":tenderId/status", handler.EditTenderStatus)
			tenders.PUT(":tenderId/rollback/:version", handler.TenderRollBack)
		}
		bids := api.Group("/bids")
		{
			bids.POST("/new", handler.CreateBid)
			bids.GET("/my", handler.GetBidsByUsername)

			bids.GET(":bidId/status", handler.GetBidStatus)
			bids.PUT(":bidId/status", handler.EditBidStatus)
			bids.PATCH(":bidId/edit", handler.EditBid)
			bids.PUT(":bidId/submit_decision", handler.SubmitDecision)
			bids.PUT(":bidId/rollback/:version", handler.BidRollBack)
		}
	}

	return router
}
func (handler *Handler) ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "server run")
}
