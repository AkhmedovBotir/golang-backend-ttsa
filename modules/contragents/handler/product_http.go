package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/contragents/service"
	"github.com/gin-gonic/gin"
)

type ContragentProductHandler struct {
	service service.ContragentProductService
}

type productStatusRequest struct {
	Status string `json:"status"`
}

func NewContragentProductHandler(s service.ContragentProductService) *ContragentProductHandler {
	return &ContragentProductHandler{service: s}
}

func (h *ContragentProductHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc) {
	grp := api.Group("/contragents/me/products")
	grp.Use(auth)
	{
		grp.POST("", h.Create)
		grp.GET("", h.GetAll)
		grp.GET("/:id", h.GetByID)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id/status", h.UpdateStatus)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *ContragentProductHandler) Create(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	var req service.ProductInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Create(contragentID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Mahsulot yaratildi", row, nil)
}

func (h *ContragentProductHandler) GetAll(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	paginated, err := h.service.GetPaginated(contragentID, page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulotlar ro'yxati olindi", gin.H{
		"items":       paginated.Items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *ContragentProductHandler) GetByID(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetByID(contragentID, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot ma'lumoti olindi", row, nil)
}

func (h *ContragentProductHandler) Update(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.ProductInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Update(contragentID, id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot yangilandi", row, nil)
}

func (h *ContragentProductHandler) UpdateStatus(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req productStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateStatus(contragentID, id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot statusi yangilandi", row, nil)
}

func (h *ContragentProductHandler) Delete(c *gin.Context) {
	contragentID, ok := getContragentID(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
		return
	}
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(contragentID, id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot o'chirildi", nil, nil)
}

func (h *ContragentProductHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrProductNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrProductStatusInvalid, service.ErrProductNameRequired, service.ErrProductDescriptionRequired,
		service.ErrProductPriceInvalid, service.ErrProductOriginalInvalid, service.ErrProductCategoryInvalid,
		service.ErrProductSubcategoryInvalid, service.ErrProductCategoryRelation, service.ErrProductQuantityInvalid,
		service.ErrProductUnitInvalid, service.ErrProductUnitSizeRequired, service.ErrProductImagesInvalid,
		service.ErrProductImageBase64Invalid, service.ErrProductKPIInvalid:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}
