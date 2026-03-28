package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type SubcategoryHandler struct {
	service service.CategoryService
}

type subcategoryStatusRequest struct {
	Status string `json:"status"`
}

func NewSubcategoryHandler(s service.CategoryService) *SubcategoryHandler {
	return &SubcategoryHandler{service: s}
}

func (h *SubcategoryHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/subcategories")
	grp.Use(auth, onlyGeneral)
	{
		grp.POST("", h.Create)
		grp.GET("", h.GetAll)
		grp.GET("/:id", h.GetByID)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id/status", h.UpdateStatus)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *SubcategoryHandler) Create(c *gin.Context) {
	var req service.SubcategoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Slug == "" || req.Status == "" || req.ParentID == 0 {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.CreateSubcategory(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Subkategoriya yaratildi", row, nil)
}

func (h *SubcategoryHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	parentRaw := c.Query("parent_id")
	var parentID *uint
	if parentRaw != "" {
		id, err := parseUintID(parentRaw)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "parent_id noto'g'ri", nil, nil)
			return
		}
		parentID = &id
	}
	paginated, err := h.service.GetSubcategories(page, limit, parentID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Subkategoriyalar ro'yxati olindi", gin.H{
		"items":       paginated.Items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *SubcategoryHandler) GetByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetSubcategoryByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Subkategoriya ma'lumoti olindi", row, nil)
}

func (h *SubcategoryHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.SubcategoryInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Name == "" || req.Slug == "" || req.Status == "" || req.ParentID == 0 {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}
	row, err := h.service.UpdateSubcategory(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Subkategoriya yangilandi", row, nil)
}

func (h *SubcategoryHandler) UpdateStatus(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req subcategoryStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}
	row, err := h.service.UpdateSubcategoryStatus(id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Subkategoriya statusi yangilandi", row, nil)
}

func (h *SubcategoryHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.DeleteSubcategory(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Subkategoriya o'chirildi", nil, nil)
}

func (h *SubcategoryHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrSubcategoryNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrInvalidStatusField, service.ErrCategoryNameRequired, service.ErrCategorySlugRequired,
		service.ErrCategorySlugExists, service.ErrCategoryParentInvalid:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}
