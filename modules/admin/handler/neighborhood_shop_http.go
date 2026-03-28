package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/domain"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type NeighborhoodShopHandler struct {
	service service.NeighborhoodShopService
}

func NewNeighborhoodShopHandler(s service.NeighborhoodShopService) *NeighborhoodShopHandler {
	return &NeighborhoodShopHandler{service: s}
}

func (h *NeighborhoodShopHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/neighborhood_shops")
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

func (h *NeighborhoodShopHandler) Create(c *gin.Context) {
	var req service.NeighborhoodShopInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Create(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Maxalla do'koni yaratildi", sanitizeNeighborhoodShop(row), nil)
}

func (h *NeighborhoodShopHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginated, err := h.service.GetPaginated(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}

	items := make([]gin.H, 0, len(paginated.Items))
	for i := range paginated.Items {
		items = append(items, sanitizeNeighborhoodShopListItem(&paginated.Items[i]))
	}

	response.JSON(c, http.StatusOK, "Maxalla do'konlari ro'yxati olindi", gin.H{
		"items":       items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *NeighborhoodShopHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Maxalla do'koni ma'lumoti olindi", sanitizeNeighborhoodShop(row), nil)
}

func (h *NeighborhoodShopHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.NeighborhoodShopInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Maxalla do'koni yangilandi", sanitizeNeighborhoodShop(row), nil)
}

func (h *NeighborhoodShopHandler) UpdateStatus(c *gin.Context) {
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
	row, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Maxalla do'koni statusi yangilandi", sanitizeNeighborhoodShop(row), nil)
}

func (h *NeighborhoodShopHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Maxalla do'koni o'chirildi", nil, nil)
}

func (h *NeighborhoodShopHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrNeighborhoodShopNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrNeighborhoodShopPhoneExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrNeighborhoodShopInvalidINN, service.ErrInvalidPhone, service.ErrNeighborhoodShopHierarchy,
		service.ErrInvalidStatusField, service.ErrNeighborhoodShopNameRequired, service.ErrNeighborhoodShopLocationIDs:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func sanitizeNeighborhoodShop(row *domain.NeighborhoodShop) gin.H {
	hasPassword := row.Password != ""
	return gin.H{
		"id":                     row.ID,
		"name":                   row.Name,
		"inn":                    row.INN,
		"region_id":              row.RegionID,
		"district_id":            row.DistrictID,
		"mfy_id":                 row.MFYID,
		"phone":                  row.Phone,
		"logo":                   row.Logo,
		"status":                 row.Status,
		"password_setup_allowed": row.PasswordSetupAllowed,
		"has_password":           hasPassword,
		"created_at":             row.CreatedAt,
		"updated_at":             row.UpdatedAt,
	}
}

func sanitizeNeighborhoodShopListItem(row *domain.NeighborhoodShop) gin.H {
	h := sanitizeNeighborhoodShop(row)
	h["logo"] = ""
	return h
}
