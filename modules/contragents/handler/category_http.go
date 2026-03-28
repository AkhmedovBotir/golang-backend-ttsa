package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/contragents/service"
	"github.com/gin-gonic/gin"
)

type ContragentCategoryHandler struct {
	service service.ContragentCategoryService
}

func NewContragentCategoryHandler(s service.ContragentCategoryService) *ContragentCategoryHandler {
	return &ContragentCategoryHandler{service: s}
}

func (h *ContragentCategoryHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc) {
	grp := api.Group("/contragents/me")
	grp.Use(auth)
	{
		grp.GET("/categories", h.GetCategories)
		grp.GET("/categories/:id", h.GetCategoryByID)
		grp.GET("/subcategories", h.GetSubcategories)
		grp.GET("/subcategories/:id", h.GetSubcategoryByID)
	}
}

func (h *ContragentCategoryHandler) GetCategories(c *gin.Context) {
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

func (h *ContragentCategoryHandler) GetCategoryByID(c *gin.Context) {
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

func (h *ContragentCategoryHandler) GetSubcategories(c *gin.Context) {
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

func (h *ContragentCategoryHandler) GetSubcategoryByID(c *gin.Context) {
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

func (h *ContragentCategoryHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrContragentCategoryNotFound, service.ErrContragentSubcategoryNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
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
