package handler

import (
	"net/http"

	"backend/internal/pkg/response"
	"backend/modules/contragents/service"
	"github.com/gin-gonic/gin"
)

type ContragentRegionDeliveryHandler struct {
	service service.ContragentRegionDeliveryService
}

func NewContragentRegionDeliveryHandler(s service.ContragentRegionDeliveryService) *ContragentRegionDeliveryHandler {
	return &ContragentRegionDeliveryHandler{service: s}
}

func (h *ContragentRegionDeliveryHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc) {
	grp := api.Group("/contragents/me")
	grp.Use(auth)
	{
		grp.GET("/regions", h.GetRegions)
		grp.GET("/districts", h.GetDistricts)
		grp.GET("/mfys", h.GetMFYs)
		grp.GET("/delivery-areas", h.GetDeliveryAreas)
		grp.PUT("/delivery-areas", h.SaveDeliveryAreas)
	}
}

func (h *ContragentRegionDeliveryHandler) GetRegions(c *gin.Context) {
	rows, err := h.service.GetRegions()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Viloyatlar ro'yxati olindi", rows, nil)
}

func (h *ContragentRegionDeliveryHandler) GetDistricts(c *gin.Context) {
	regionRaw := c.Query("region_id")
	var regionID *uint
	if regionRaw != "" {
		id, err := parseUintID(regionRaw)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "region_id noto'g'ri", nil, nil)
			return
		}
		regionID = &id
	}
	rows, err := h.service.GetDistricts(regionID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Tumanlar ro'yxati olindi", rows, nil)
}

func (h *ContragentRegionDeliveryHandler) GetMFYs(c *gin.Context) {
	districtRaw := c.Query("district_id")
	var districtID *uint
	if districtRaw != "" {
		id, err := parseUintID(districtRaw)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "district_id noto'g'ri", nil, nil)
			return
		}
		districtID = &id
	}
	rows, err := h.service.GetMFYs(districtID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "MFYlar ro'yxati olindi", rows, nil)
}

func (h *ContragentRegionDeliveryHandler) GetDeliveryAreas(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	result, err := h.service.GetDeliveryAreas(contragentID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Yetkazib berish hududlari olindi", result, nil)
}

func (h *ContragentRegionDeliveryHandler) SaveDeliveryAreas(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	var req service.DeliveryAreasInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	result, err := h.service.SaveDeliveryAreas(contragentID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Yetkazib berish hududlari saqlandi", result, nil)
}

func (h *ContragentRegionDeliveryHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrRegionIDsRequired, service.ErrDistrictIDsRequired, service.ErrDeliveryHierarchyInvalid:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	case service.ErrDeliveryRegionNotFound, service.ErrDeliveryDistrictNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func getContragentID(c *gin.Context) (uint, bool) {
	idVal, exists := c.Get("contragent_id")
	if !exists {
		return 0, false
	}
	id, ok := idVal.(uint)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}
