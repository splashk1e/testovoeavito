package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"testovoe.com/internal/models"
)

func (handler *Handler) CreateTender(ctx *gin.Context) {
	username := ctx.Param("username")
	var input models.Tender
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	tender, err := handler.services.Tenders.CreateTender(input, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}

func (handler *Handler) GetTenders(ctx *gin.Context) {
	var start int
	var limit int
	var serviceTypes []string
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
	for i := 0; i < 3; i++ {
		serviceType := ctx.Param("service_type")
		serviceTypes = append(serviceTypes, serviceType)
	}
	tender, err := handler.services.Tenders.GetTenders(serviceTypes)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	output := tender[start : start+limit]
	ctx.JSON(http.StatusOK, output)
}

func (handler *Handler) GetTendersByUsername(ctx *gin.Context) {
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
	tender, err := handler.services.Tenders.GetTendersByUsername(username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	output := tender[start : start+limit]
	ctx.JSON(http.StatusOK, output)
}
func (handler *Handler) GetTenderStatus(ctx *gin.Context) {

	username := ctx.Param("username")
	tenderId := ctx.Param("tenderId")

	tender, err := handler.services.Tenders.GetTenderById(tenderId, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender.Status)
}
func (handler *Handler) EditTenderStatus(ctx *gin.Context) {
	username := ctx.Param("username")
	tenderId := ctx.Param("tenderId")
	status := ctx.Param("status")
	tender, err := handler.services.Tenders.EditTenderStatusById(tenderId, username, status)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}

type changes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}

func (handler *Handler) EditTender(ctx *gin.Context) {
	id := ctx.Param("tenderId")
	username := ctx.Param("username")
	var input changes
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	changesmap := map[string]string{
		"name":         input.Name,
		"decription":   input.Description,
		"service_type": input.ServiceType,
	}
	tender, err := handler.services.Tenders.EditTenderById(id, changesmap, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}
func (handler *Handler) TenderRollBack(ctx *gin.Context) {
	username := ctx.Param("username")
	tenderId := ctx.Param("tenderId")
	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	tender, err := handler.services.Tenders.TenderRollBack(tenderId, version, username)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, tender)
}
