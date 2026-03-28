package handler

import (
	"net/http"
	"strconv"

	"backend/internal/pkg/response"
	"backend/modules/admin/domain"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	service service.AgentService
}

func NewAgentHandler(s service.AgentService) *AgentHandler {
	return &AgentHandler{service: s}
}

func (h *AgentHandler) RegisterRoutes(api *gin.RouterGroup, auth gin.HandlerFunc, onlyGeneral gin.HandlerFunc) {
	grp := api.Group("/agents")
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

func (h *AgentHandler) Create(c *gin.Context) {
	var req service.AgentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Create(req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "Agent yaratildi", sanitizeAgent(row), nil)
}

func (h *AgentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginated, err := h.service.GetPaginated(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}

	items := make([]gin.H, 0, len(paginated.Items))
	for i := range paginated.Items {
		items = append(items, sanitizeAgent(&paginated.Items[i]))
	}

	response.JSON(c, http.StatusOK, "Agentlar ro'yxati olindi", gin.H{
		"items":       items,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *AgentHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Agent ma'lumoti olindi", sanitizeAgent(row), nil)
}

func (h *AgentHandler) Update(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	var req service.AgentInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	row, err := h.service.Update(id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Agent yangilandi", sanitizeAgent(row), nil)
}

func (h *AgentHandler) UpdateStatus(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, "Agent statusi yangilandi", sanitizeAgent(row), nil)
}

func (h *AgentHandler) Delete(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}
	if err = h.service.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Agent o'chirildi", nil, nil)
}

func (h *AgentHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrAgentNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case service.ErrAgentPhoneExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrInvalidPhone, service.ErrAgentHierarchy, service.ErrInvalidStatusField,
		service.ErrAgentNameRequired, service.ErrAgentLocationIDs:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func sanitizeAgent(row *domain.Agent) gin.H {
	hasPassword := row.Password != ""
	return gin.H{
		"id":                     row.ID,
		"name":                   row.Name,
		"viloyat_id":             row.RegionID,
		"tuman_id":               row.DistrictID,
		"mfy_id":                 row.MFYID,
		"phone":                  row.Phone,
		"status":                 row.Status,
		"password_setup_allowed": row.PasswordSetupAllowed,
		"has_password":           hasPassword,
		"created_at":             row.CreatedAt,
		"updated_at":             row.UpdatedAt,
	}
}
