package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/pkg/response"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type MarketplaceSearchHandler struct {
	searchSvc     service.MarketplaceSearchService
	contragentSvc service.MarketplaceContragentBrowseService
}

func NewMarketplaceSearchHandler(
	searchSvc service.MarketplaceSearchService,
	contragentSvc service.MarketplaceContragentBrowseService,
) *MarketplaceSearchHandler {
	return &MarketplaceSearchHandler{searchSvc: searchSvc, contragentSvc: contragentSvc}
}

func (h *MarketplaceSearchHandler) RegisterRoutes(api *gin.RouterGroup) {
	grp := api.Group("/marketplace")
	{
		grp.GET("/search", h.UnifiedSearch)
		grp.GET("/contragents", h.ListContragents)
	}
}

func (h *MarketplaceSearchHandler) UnifiedSearch(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	limitPerType, _ := strconv.Atoi(c.DefaultQuery("limit_per_type", "10"))
	mask := service.ParseSearchTypes(c.Query("types"))

	out, err := h.searchSvc.UnifiedSearch(q, limitPerType, mask)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Qidiruv natijalari", out, nil)
}

func (h *MarketplaceSearchHandler) ListContragents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	nestedLimit, _ := strconv.Atoi(c.DefaultQuery("nested_limit", "30"))
	inc := service.ParseContragentInclude(c.Query("include"))

	var query *string
	if v := strings.TrimSpace(c.Query("q")); v != "" {
		query = &v
	}

	out, err := h.contragentSvc.List(query, page, limit, nestedLimit, inc)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Kontragentlar ro'yxati olindi", gin.H{
		"items":       out.Items,
		"total":       out.Total,
		"page":        out.Page,
		"limit":       out.Limit,
		"total_pages": out.TotalPages,
	}, nil)
}
