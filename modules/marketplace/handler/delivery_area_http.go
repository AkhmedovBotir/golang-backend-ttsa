package handler

import (
	"errors"
	"net/http"

	"backend/internal/pkg/response"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type DeliveryAreaHandler struct {
	service service.DeliveryAreaService
}

func NewDeliveryAreaHandler(service service.DeliveryAreaService) *DeliveryAreaHandler {
	return &DeliveryAreaHandler{service: service}
}

func (h *DeliveryAreaHandler) RegisterRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	grp := api.Group("/marketplace/me/delivery-areas")
	grp.Use(authMiddleware)
	{
		grp.POST("", h.Create)
		grp.GET("", h.List)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id/set-default", h.SetDefault)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *DeliveryAreaHandler) Create(c *gin.Context) {
	userID, ok := getMarketplaceUserID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	var req service.DeliveryAreaInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.Create(userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "yetkazib berish hududi saqlandi", out, nil)
}

func (h *DeliveryAreaHandler) List(c *gin.Context) {
	userID, ok := getMarketplaceUserID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	rows, err := h.service.List(userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "yetkazib berish hududlari olindi", rows, nil)
}

func (h *DeliveryAreaHandler) Update(c *gin.Context) {
	userID, ok := getMarketplaceUserID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "id noto'g'ri", nil, nil)
		return
	}
	var req service.DeliveryAreaInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.Update(userID, id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "yetkazib berish hududi yangilandi", out, nil)
}

func (h *DeliveryAreaHandler) Delete(c *gin.Context) {
	userID, ok := getMarketplaceUserID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "id noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(userID, id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "yetkazib berish hududi o'chirildi", gin.H{"deleted": true}, nil)
}

func (h *DeliveryAreaHandler) SetDefault(c *gin.Context) {
	userID, ok := getMarketplaceUserID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "id noto'g'ri", nil, nil)
		return
	}
	out, err := h.service.SetDefault(userID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "asosiy manzil belgilandi", out, nil)
}

func (h *DeliveryAreaHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrDeliveryAreaNameRequired),
		errors.Is(err, service.ErrDeliveryLocationIDs),
		errors.Is(err, service.ErrDeliveryHierarchy):
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	case errors.Is(err, service.ErrDeliveryAreaNotFound):
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func getMarketplaceUserID(c *gin.Context) (uint, bool) {
	idVal, ok := c.Get("marketplace_user_id")
	if !ok {
		return 0, false
	}
	id, ok := idVal.(uint)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}
