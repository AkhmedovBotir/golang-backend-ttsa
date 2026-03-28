package handler

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/pkg/response"
	"backend/internal/pkg/security"
	"backend/modules/admin/domain"
	"backend/modules/admin/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service        service.AdminService
	jwtSecret      string
	jwtExpireHours int
}

func NewAdminHandler(s service.AdminService, jwtSecret string, jwtExpireHours int) *AdminHandler {
	return &AdminHandler{service: s, jwtSecret: jwtSecret, jwtExpireHours: jwtExpireHours}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type statusRequest struct {
	Status string `json:"status"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	admin, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredential:
			response.JSON(c, http.StatusUnauthorized, err.Error(), nil, nil)
		case service.ErrInactiveAdmin:
			response.JSON(c, http.StatusForbidden, err.Error(), nil, nil)
		default:
			response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		}
		return
	}

	token, err := security.GenerateToken(h.jwtSecret, admin.ID, admin.Role, h.jwtExpireHours)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Token yaratishda xatolik", nil, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Muvaffaqiyatli login qilindi", gin.H{
		"token": token,
		"admin": sanitizeAdmin(admin),
	}, nil)
}

func (h *AdminHandler) Create(c *gin.Context) {
	var req service.CreateAdminInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	if req.Name == "" || req.Role == "" || req.Phone == "" || req.Username == "" || req.Password == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}

	admin, err := h.service.Create(req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusCreated, "Admin muvaffaqiyatli yaratildi", sanitizeAdmin(admin), nil)
}

func (h *AdminHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginated, err := h.service.GetPaginated(page, limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
		return
	}

	result := make([]gin.H, 0, len(paginated.Items))
	for _, a := range paginated.Items {
		admin := a
		result = append(result, sanitizeAdmin(&admin))
	}

	response.JSON(c, http.StatusOK, "Adminlar ro'yxati olindi", gin.H{
		"items":       result,
		"total":       paginated.Total,
		"page":        paginated.Page,
		"limit":       paginated.Limit,
		"total_pages": paginated.TotalPages,
	}, nil)
}

func (h *AdminHandler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	admin, err := h.service.GetByID(id)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Admin ma'lumoti olindi", sanitizeAdmin(admin), nil)
}

func (h *AdminHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	var req service.UpdateAdminInput
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}

	if req.Name == "" || req.Role == "" || req.Phone == "" || req.Username == "" || req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "Majburiy maydonlar to'ldirilmagan", nil, nil)
		return
	}

	admin, err := h.service.Update(id, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Admin muvaffaqiyatli yangilandi", sanitizeAdmin(admin), nil)
}

func (h *AdminHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	if err = h.service.Delete(id); err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Admin muvaffaqiyatli o'chirildi", nil, nil)
}

func (h *AdminHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.JSON(c, http.StatusBadRequest, "ID noto'g'ri", nil, nil)
		return
	}

	var req statusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "So'rov formati noto'g'ri", nil, err.Error())
		return
	}
	if req.Status == "" {
		response.JSON(c, http.StatusBadRequest, "status majburiy", nil, nil)
		return
	}

	admin, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "Admin statusi yangilandi", sanitizeAdmin(admin), nil)
}

func (h *AdminHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.JSON(c, http.StatusUnauthorized, "Authorization header topilmadi", nil, nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.JSON(c, http.StatusUnauthorized, "Authorization formati noto'g'ri", nil, nil)
			c.Abort()
			return
		}

		claims, err := security.ParseToken(h.jwtSecret, parts[1])
		if err != nil {
			response.JSON(c, http.StatusUnauthorized, "Token yaroqsiz", nil, nil)
			c.Abort()
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("admin_role", claims.Role)
		c.Next()
	}
}

func GeneralOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("admin_role")
		if !ok || role != domain.RoleGeneral {
			response.JSON(c, http.StatusForbidden, "Faqat general admin ushbu amalni bajara oladi", nil, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *AdminHandler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidPhone, service.ErrInvalidRole, service.ErrInvalidStatus:
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	case service.ErrPhoneExists, service.ErrUsernameExists:
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case service.ErrNotFound:
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, "Serverda xatolik yuz berdi", nil, err.Error())
	}
}

func parseID(raw string) (uint, error) {
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func sanitizeAdmin(admin *domain.Admin) gin.H {
	return gin.H{
		"id":         admin.ID,
		"name":       admin.Name,
		"role":       admin.Role,
		"phone":      admin.Phone,
		"username":   admin.Username,
		"status":     admin.Status,
		"created_at": admin.CreatedAt,
		"updated_at": admin.UpdatedAt,
	}
}
