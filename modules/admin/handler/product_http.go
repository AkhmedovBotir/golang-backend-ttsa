package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/pkg/response"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.AdminProductService
}

type productStatusRequest struct {
	Status string `json:"status"`
}

type rejectProductRequest struct {
	RejectionReason string `json:"rejection_reason"`
}

func NewProductHandler(s service.AdminProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/products")
	grp.Use(auth, onlyGeneral)
	{
		grp.POST("", h.Create)
		grp.GET("", h.GetAll)
		grp.GET("/:id", h.GetByID)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id/status", h.UpdateStatus)
		grp.PATCH("/:id/approve", h.Approve)
		grp.PATCH("/:id/reject", h.Reject)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req service.AdminProductInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Create(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Mahsulot yaratildi", row, nil)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var contragentID *uint
	contrRaw := strings.TrimSpace(c.Query("contragent_id"))
	if contrRaw != "" {
		id, err := parseUintID(contrRaw)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "contragent_id noto'g'ri", nil, nil)
			return
		}
		contragentID = &id
	}

	var moderationStatus *string
	modRaw := strings.TrimSpace(c.Query("moderation_status"))
	if modRaw != "" {
		moderationStatus = &modRaw
	}

	paginated, err := h.service.GetPaginated(page, limit, contragentID, moderationStatus)
	if err != nil {
		h.handleError(c, err)
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

func (h *ProductHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Mahsulot ma'lumoti olindi", row, nil)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.AdminProductInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot yangilandi", row, nil)
}

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
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
	row, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot statusi yangilandi", row, nil)
}

func (h *ProductHandler) Approve(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.Approve(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot tasdiqlandi", row, nil)
}

func (h *ProductHandler) Reject(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req rejectProductRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	row, err := h.service.Reject(id, req.RejectionReason)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot rad etildi", row, nil)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot o'chirildi", nil, nil)
}

func (h *ProductHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrAdminProductNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrAdminProductNameRequired, service.ErrAdminProductDescriptionRequired,
		service.ErrAdminProductPriceInvalid, service.ErrAdminProductOriginalInvalid, service.ErrAdminProductCategoryInvalid,
		service.ErrAdminProductSubcategoryInvalid, service.ErrAdminProductCategoryRelation, service.ErrAdminProductQuantityInvalid,
		service.ErrAdminProductUnitInvalid, service.ErrAdminProductUnitSizeRequired, service.ErrAdminProductImagesInvalid,
		service.ErrAdminProductImageBase64Invalid, service.ErrAdminProductKPIInvalid, service.ErrAdminProductStatusInvalid,
		service.ErrAdminProductContragentInvalid, service.ErrAdminProductModerationInvalid, service.ErrAdminRejectReasonRequired:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}
