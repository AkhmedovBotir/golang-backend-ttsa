package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	service service.RegionService
}

func NewRegionHandler(service service.RegionService) *RegionHandler {
	return &RegionHandler{service: service}
}

func (h *RegionHandler) RegisterRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	_ = authMiddleware
	grp := api.Group("/marketplace")
	{
		grp.GET("/regions", h.GetRegions)
		grp.GET("/districts", h.GetDistricts)
		grp.GET("/mfys", h.GetMFYs)
	}
}

func (h *RegionHandler) GetRegions(c *gin.Context) {
	rows, err := h.service.GetRegions()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Viloyatlar ro'yxati olindi", rows, nil)
}

func (h *RegionHandler) GetDistricts(c *gin.Context) {
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

func (h *RegionHandler) GetMFYs(c *gin.Context) {
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

func parseUintID(raw string) (uint, error) {
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
