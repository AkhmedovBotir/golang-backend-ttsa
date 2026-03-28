package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/pkg/response"
	adminDomain "backend/modules/admin/domain"
	"backend/modules/marketplace/repository"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type MarketplaceProductHandler struct {
	service service.MarketplaceProductService
}

func NewMarketplaceProductHandler(s service.MarketplaceProductService) *MarketplaceProductHandler {
	return &MarketplaceProductHandler{service: s}
}

func (h *MarketplaceProductHandler) RegisterRoutes(api *gin.RouterGroup) {
	grp := api.Group("/marketplace")
	{
		grp.GET("/products", h.List)
		grp.GET("/products/:id", h.GetByID)
	}
}

func (h *MarketplaceProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := repository.MarketplaceProductFilter{}

	if v := strings.TrimSpace(c.Query("category_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "category_id noto'g'ri", nil, nil)
			return
		}
		filter.CategoryID = &id
	}
	if v := strings.TrimSpace(c.Query("subcategory_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "subcategory_id noto'g'ri", nil, nil)
			return
		}
		filter.SubcategoryID = &id
	}
	if v := strings.TrimSpace(c.Query("contragent_id")); v != "" {
		id, err := parseUintID(v)
		if err != nil {
			response.JSON(c, http.StatusBadRequest, "contragent_id noto'g'ri", nil, nil)
			return
		}
		filter.ContragentID = &id
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		if status != adminDomain.StatusActive && status != adminDomain.StatusInactive {
			response.JSON(c, http.StatusBadRequest, "status noto'g'ri", nil, nil)
			return
		}
		filter.Status = &status
	}

	if q := strings.TrimSpace(c.Query("q")); q != "" {
		filter.Query = &q
	}

	out, err := h.service.ListApproved(service.MarketplaceProductFilter(filter), page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulotlar ro'yxati olindi", gin.H{
		"items":       out.Items,
		"total":       out.Total,
		"page":        out.Page,
		"limit":       out.Limit,
		"total_pages": out.TotalPages,
	}, nil)
}

func (h *MarketplaceProductHandler) GetByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "id noto'g'ri", nil, nil)
		return
	}

	out, err := h.service.GetApprovedByID(id)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	if out == nil {
		response.JSON(c, http.StatusNotFound, "Mahsulot topilmadi", nil, nil)
		return
	}
	response.JSON(c, http.StatusOK, "Mahsulot ma'lumoti olindi", out, nil)
}
