package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type MarketplaceCategoryHandler struct {
	service service.MarketplaceCategoryService
}

func NewMarketplaceCategoryHandler(s service.MarketplaceCategoryService) *MarketplaceCategoryHandler {
	return &MarketplaceCategoryHandler{service: s}
}

func (h *MarketplaceCategoryHandler) RegisterRoutes(api *gin.RouterGroup) {
	grp := api.Group("/marketplace")
	{
		grp.GET("/categories", h.GetCategories)
		grp.GET("/categories/:id", h.GetCategoryByID)
		grp.GET("/subcategories", h.GetSubcategories)
		grp.GET("/subcategories/:id", h.GetSubcategoryByID)
	}
}

func (h *MarketplaceCategoryHandler) GetCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	paginated, err := h.service.GetCategories(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Kategoriyalar ro'yxati olindi", gin.H{
		"items":       paginated.Items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *MarketplaceCategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	row, err := h.service.GetCategoryByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Kategoriya ma'lumoti olindi", row, nil)
}

func (h *MarketplaceCategoryHandler) GetSubcategories(c *gin.Context) {
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

func (h *MarketplaceCategoryHandler) GetSubcategoryByID(c *gin.Context) {
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

func (h *MarketplaceCategoryHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrMarketplaceCategoryNotFound, service.ErrMarketplaceSubcategoryNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}
