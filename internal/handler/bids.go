package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"testovoe.com/internal/models"
)

func (handler *Handler) CreateBid(ctx *gin.Context) {
	username := ctx.Param("username")
	var input models.Bid
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	bid, err := handler.services.Bids.CreateBid(input, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, bid)
}
func (handler *Handler) GetBidsByTenderId(ctx *gin.Context) {
	var start int
	var limit int
	var username string
	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	start, err = strconv.Atoi(ctx.Param("offset"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	username = ctx.Param("username")
	tenderId := ctx.Param("tenderId")
	bids, err := handler.services.Bids.GetBidByTenderId(tenderId, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	output := bids[start : start+limit]
	ctx.JSON(http.StatusOK, output)
}
func (handler *Handler) GetBidsByUsername(ctx *gin.Context) {
	var start int
	var limit int
	var username string
	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	start, err = strconv.Atoi(ctx.Param("offset"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	username = ctx.Param("username")
	bid, err := handler.services.Bids.GetBidsByUsername(username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	output := bid[start : start+limit]
	ctx.JSON(http.StatusOK, output)
}
func (handler *Handler) GetBidStatus(ctx *gin.Context) {
	username := ctx.Param("username")
	id := ctx.Param("bidId")

	bid, err := handler.services.Bids.GetBidById(id, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, bid.Status)
}
func (handler *Handler) EditBidStatus(ctx *gin.Context) {
	username := ctx.Param("username")
	bidId := ctx.Param("bidId")
	status := ctx.Param("status")
	var input models.Tender
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	bid, err := handler.services.Bids.ChangeBidStatus(bidId, status, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, bid)
}
func (handler *Handler) EditBid(ctx *gin.Context) {
	id := ctx.Param("bidId")
	username := ctx.Param("username")
	var input changes
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	changesmap := map[string]string{
		"name":       input.Name,
		"decription": input.Description,
	}
	tender, err := handler.services.Bids.EditBidById(id, changesmap, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}
func (handler *Handler) BidRollBack(ctx *gin.Context) {
	username := ctx.Param("username")
	bidId := ctx.Param("bidId")
	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	tender, err := handler.services.Bids.BidRollBack(bidId, version, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}
func (handler *Handler) SubmitDecision(ctx *gin.Context) {

}
