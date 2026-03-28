package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/pkg/response"
	"backend/modules/admin/repository"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type MarketplaceUserHandler struct {
	service service.MarketplaceUserService
}

type marketplaceUserStatusRequest struct {
	Status string `json:"status"`
}

func NewMarketplaceUserHandler(s service.MarketplaceUserService) *MarketplaceUserHandler {
	return &MarketplaceUserHandler{service: s}
}

func (h *MarketplaceUserHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/marketplace-users")
	grp.Use(auth, onlyGeneral)
	{
		grp.GET("", h.GetAll)
		grp.GET("/:id", h.GetByID)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id/status", h.UpdateStatus)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *MarketplaceUserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := repository.MarketplaceUserFilter{}
	if v := strings.TrimSpace(c.Query("status")); v != "" {
		filter.Status = &v
	}
	if v := strings.TrimSpace(c.Query("region_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "region_id noto'g'ri", nil, nil)
			return
		}
		filter.RegionID = &id
	}
	if v := strings.TrimSpace(c.Query("district_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "district_id noto'g'ri", nil, nil)
			return
		}
		filter.DistrictID = &id
	}
	if v := strings.TrimSpace(c.Query("mfy_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "mfy_id noto'g'ri", nil, nil)
			return
		}
		filter.MFYID = &id
	}
	if v := strings.TrimSpace(c.Query("phone")); v != "" {
		filter.Phone = &v
	}
	if v := strings.TrimSpace(c.Query("q")); v != "" {
		filter.Query = &v
	}

	out, err := h.service.GetPaginated(page, limit, filter)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Marketplace userlar ro'yxati olindi", gin.H{
		"items":       out.Items,
		"total":       out.Total,
		"page":        out.Page,
		"limit":       out.Limit,
		"total_pages": out.TotalPages,
	}, nil)
}

func (h *MarketplaceUserHandler) GetByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Marketplace user ma'lumoti olindi", row, nil)
}

func (h *MarketplaceUserHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.MarketplaceUserInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Marketplace user yangilandi", row, nil)
}

func (h *MarketplaceUserHandler) UpdateStatus(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req marketplaceUserStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if strings.TrimSpace(req.Status) == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Marketplace user statusi yangilandi", row, nil)
}

func (h *MarketplaceUserHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Marketplace user o'chirildi", nil, nil)
}

func (h *MarketplaceUserHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrMarketplaceUserNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrMarketplaceUserPhoneExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrMarketplaceUserNameRequired, service.ErrMarketplaceUserGender,
		service.ErrMarketplaceUserBirthDate, service.ErrMarketplaceUserLocation,
		service.ErrInvalidPhone, service.ErrInvalidStatusField:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}
