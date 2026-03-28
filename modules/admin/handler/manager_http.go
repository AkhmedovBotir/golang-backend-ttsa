package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/domain"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type ManagerHandler struct {
	service service.ManagerService
}

func NewManagerHandler(s service.ManagerService) *ManagerHandler {
	return &ManagerHandler{service: s}
}

func (h *ManagerHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/managers")
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

func (h *ManagerHandler) Create(c *gin.Context) {
	var req service.ManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Create(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Menejer yaratildi", sanitizeManager(row), nil)
}

func (h *ManagerHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginated, err := h.service.GetPaginated(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}

	items := make([]gin.H, 0, len(paginated.Items))
	for i := range paginated.Items {
		items = append(items, sanitizeManager(&paginated.Items[i]))
	}

	response.JSON(c, http.StatusOK, "Menejerlar ro'yxati olindi", gin.H{
		"items":       items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *ManagerHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Menejer ma'lumoti olindi", sanitizeManager(row), nil)
}

func (h *ManagerHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.ManagerInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Menejer yangilandi", sanitizeManager(row), nil)
}

func (h *ManagerHandler) UpdateStatus(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Menejer statusi yangilandi", sanitizeManager(row), nil)
}

func (h *ManagerHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Menejer o'chirildi", nil, nil)
}

func (h *ManagerHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrManagerNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrManagerPhoneExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrInvalidPhone, service.ErrManagerHierarchy, service.ErrInvalidStatusField,
		service.ErrManagerNameRequired, service.ErrManagerRegionID:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func sanitizeManager(row *domain.Manager) gin.H {
	hasPassword := row.Password != ""
	return gin.H{
		"id":                     row.ID,
		"name":                   row.Name,
		"viloyat_id":             row.RegionID,
		"phone":                  row.Phone,
		"status":                 row.Status,
		"password_setup_allowed": row.PasswordSetupAllowed,
		"has_password":           hasPassword,
		"created_at":             row.CreatedAt,
		"updated_at":             row.UpdatedAt,
	}
}
