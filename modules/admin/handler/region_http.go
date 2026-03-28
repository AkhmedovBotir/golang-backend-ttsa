package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	service service.RegionService
}

func NewRegionHandler(s service.RegionService) *RegionHandler {
	return &RegionHandler{service: s}
}

func (h *RegionHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("")
	grp.Use(auth, onlyGeneral)
	{
		grp.POST("/regions", h.CreateRegion)
		grp.GET("/regions", h.GetRegions)
		grp.GET("/regions/:id", h.GetRegionByID)
		grp.PUT("/regions/:id", h.UpdateRegion)
		grp.PATCH("/regions/:id/status", h.UpdateRegionStatus)
		grp.DELETE("/regions/:id", h.DeleteRegion)

		grp.POST("/districts", h.CreateDistrict)
		grp.GET("/districts", h.GetDistricts)
		grp.GET("/districts/:id", h.GetDistrictByID)
		grp.PUT("/districts/:id", h.UpdateDistrict)
		grp.PATCH("/districts/:id/status", h.UpdateDistrictStatus)
		grp.DELETE("/districts/:id", h.DeleteDistrict)

		grp.POST("/mfys", h.CreateMFY)
		grp.GET("/mfys", h.GetMFYs)
		grp.GET("/mfys/:id", h.GetMFYByID)
		grp.PUT("/mfys/:id", h.UpdateMFY)
		grp.PATCH("/mfys/:id/status", h.UpdateMFYStatus)
		grp.DELETE("/mfys/:id", h.DeleteMFY)
	}
}

type regionStatusRequest struct {
	Status string `json:"status"`
}

func (h *RegionHandler) CreateRegion(c *gin.Context) {
	var req service.RegionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.CreateRegion(req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Viloyat yaratildi", row, nil)
}

func (h *RegionHandler) GetRegions(c *gin.Context) {
	rows, err := h.service.GetRegions()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Viloyatlar ro'yxati olindi", rows, nil)
}

func (h *RegionHandler) GetRegionByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetRegionByID(id)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Viloyat ma'lumoti olindi", row, nil)
}

func (h *RegionHandler) UpdateRegion(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.RegionInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.UpdateRegion(id, req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Viloyat yangilandi", row, nil)
}

func (h *RegionHandler) DeleteRegion(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.DeleteRegion(id); err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Viloyat o'chirildi", nil, nil)
}

func (h *RegionHandler) UpdateRegionStatus(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	var req regionStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateRegionStatus(id, req.Status)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Viloyat statusi yangilandi", row, nil)
}

func (h *RegionHandler) CreateDistrict(c *gin.Context) {
	var req service.DistrictInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.CreateDistrict(req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Tuman yaratildi", row, nil)
}

func (h *RegionHandler) GetDistricts(c *gin.Context) {
	rows, err := h.service.GetDistricts()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Tumanlar ro'yxati olindi", rows, nil)
}

func (h *RegionHandler) GetDistrictByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetDistrictByID(id)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Tuman ma'lumoti olindi", row, nil)
}

func (h *RegionHandler) UpdateDistrict(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.DistrictInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.UpdateDistrict(id, req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Tuman yangilandi", row, nil)
}

func (h *RegionHandler) DeleteDistrict(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.DeleteDistrict(id); err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Tuman o'chirildi", nil, nil)
}

func (h *RegionHandler) UpdateDistrictStatus(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	var req regionStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateDistrictStatus(id, req.Status)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Tuman statusi yangilandi", row, nil)
}

func (h *RegionHandler) CreateMFY(c *gin.Context) {
	var req service.MFYInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.CreateMFY(req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "MFY yaratildi", row, nil)
}

func (h *RegionHandler) GetMFYs(c *gin.Context) {
	rows, err := h.service.GetMFYs()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "MFYlar ro'yxati olindi", rows, nil)
}

func (h *RegionHandler) GetMFYByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetMFYByID(id)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "MFY ma'lumoti olindi", row, nil)
}

func (h *RegionHandler) UpdateMFY(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.MFYInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Code == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.UpdateMFY(id, req)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "MFY yangilandi", row, nil)
}

func (h *RegionHandler) DeleteMFY(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.DeleteMFY(id); err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "MFY o'chirildi", nil, nil)
}

func (h *RegionHandler) UpdateMFYStatus(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	var req regionStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateMFYStatus(id, req.Status)
	if err != nil {
		h.handleRegionError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "MFY statusi yangilandi", row, nil)
}

func (h *RegionHandler) handleRegionError(c *gin.Context, err error) {
	switch err {
	case service.ErrRegionNotFound, service.ErrDistrictNotFound, service.ErrMFYNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrInvalidStatusField, service.ErrRegionIDRequired, service.ErrDistrictIDRequired:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func parseUintID(raw string) (uint, error) {
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
