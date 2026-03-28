package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/domain"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type ContragentHandler struct {
	service service.ContragentService
}

func NewContragentHandler(s service.ContragentService) *ContragentHandler {
	return &ContragentHandler{service: s}
}

func (h *ContragentHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/contragents")
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

func (h *ContragentHandler) Create(c *gin.Context) {
	var req service.ContragentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Create(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Kontragent yaratildi", sanitizeContragent(row), nil)
}

func (h *ContragentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginated, err := h.service.GetPaginated(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}

	items := make([]gin.H, 0, len(paginated.Items))
	for i := range paginated.Items {
		items = append(items, sanitizeContragentListItem(&paginated.Items[i]))
	}

	response.JSON(c, http.StatusOK, "Kontragentlar ro'yxati olindi", gin.H{
		"items":       items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *ContragentHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Kontragent ma'lumoti olindi", sanitizeContragent(row), nil)
}

func (h *ContragentHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.ContragentInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Kontragent yangilandi", sanitizeContragent(row), nil)
}

func (h *ContragentHandler) UpdateStatus(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Kontragent statusi yangilandi", sanitizeContragent(row), nil)
}

func (h *ContragentHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Kontragent o'chirildi", nil, nil)
}

func (h *ContragentHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrContragentNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrContragentPhoneExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrInvalidINN, service.ErrInvalidPhone, service.ErrContragentHierarchy, service.ErrInvalidStatusField,
		service.ErrContragentNameRequired, service.ErrContragentActivityTypeID, service.ErrContragentLocationIDs,
		service.ErrContragentActivityTypeNotFound:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func sanitizeContragent(row *domain.Contragent) gin.H {
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
		"activity_type_id":       row.ActivityTypeID,
		"status":                 row.Status,
		"password_setup_allowed": row.PasswordSetupAllowed,
		"has_password":           hasPassword,
		"created_at":             row.CreatedAt,
		"updated_at":             row.UpdatedAt,
	}
}

func sanitizeContragentListItem(row *domain.Contragent) gin.H {
	h := sanitizeContragent(row)
	h["logo"] = ""
	return h
}
