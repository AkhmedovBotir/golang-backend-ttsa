package handler

import (
	"errors"
	"net/http"
	"strings"

	"backend/internal/pkg/response"
	"backend/modules/marketplace/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	pub := api.Group("/marketplace/auth")
	{
		pub.POST("/entry", h.EntryStart)
		pub.POST("/entry/verify", h.EntryVerify)
		pub.POST("/phone/check", h.CheckPhone)
		pub.POST("/send-code", h.SendCode)
		pub.POST("/verify-code", h.VerifyCode)
		pub.POST("/resend-code", h.ResendCode)
		pub.POST("/login", h.LoginBySMS)
		pub.POST("/register", h.Register)
	}
	me := api.Group("/marketplace/me")
	me.Use(authMiddleware)
	{
		me.GET("/profile", h.GetMe)
		me.PUT("/profile", h.UpdateMe)
		me.GET("/avatar", h.GetMyAvatar)
		me.PUT("/avatar", h.UpdateMyAvatar)
		me.DELETE("/avatar", h.DeleteMyAvatar)
	}
}

type phoneRequest struct {
	Phone string `json:"phone"`
}

type sendCodeRequest struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
}

type verifyCodeRequest struct {
	Phone   string `json:"phone"`
	Code    string `json:"code"`
	Purpose string `json:"purpose"`
}

type registerRequest struct {
	Phone      string `json:"phone"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	BirthDate  string `json:"birth_date"`
}

type entryVerifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type updateMeRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	BirthDate  string `json:"birth_date"`
}

type avatarRequest struct {
	Avatar string `json:"avatar"`
}

func (h *AuthHandler) EntryStart(c *gin.Context) {
	var req phoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.EntryStart(req.Phone)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "sms kod yuborildi", out, nil)
}

func (h *AuthHandler) EntryVerify(c *gin.Context) {
	var req entryVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.EntryVerify(req.Phone, req.Code)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if out.NeedsRegistration {
		response.JSON(c, http.StatusOK, "kod tasdiqlandi, registratsiyani yakunlang", out, nil)
		return
	}
	response.JSON(c, http.StatusOK, "muvaffaqiyatli login", out, nil)
}

func (h *AuthHandler) CheckPhone(c *gin.Context) {
	var req phoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	result, err := h.service.CheckPhone(req.Phone)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "ok", result, nil)
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	err := h.service.SendCode(service.SendCodeInput{
		Phone:   req.Phone,
		Purpose: req.Purpose,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "sms kod yuborildi", gin.H{"sent": true}, nil)
}

func (h *AuthHandler) ResendCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	err := h.service.ResendCode(service.SendCodeInput{
		Phone:   req.Phone,
		Purpose: req.Purpose,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "sms kod qayta yuborildi", gin.H{"sent": true}, nil)
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req verifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	err := h.service.VerifyCode(service.VerifyCodeInput{
		Phone:   req.Phone,
		Code:    req.Code,
		Purpose: req.Purpose,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "kod tasdiqlandi", gin.H{"verified": true}, nil)
}

func (h *AuthHandler) LoginBySMS(c *gin.Context) {
	var req phoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.LoginBySMS(service.LoginInput{Phone: req.Phone})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "muvaffaqiyatli login", out, nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.Register(service.RegisterInput{
		Phone:      req.Phone,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Gender:     req.Gender,
		RegionID:   req.RegionID,
		DistrictID: req.DistrictID,
		MFYID:      req.MFYID,
		BirthDate:  req.BirthDate,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, "ro'yxatdan o'tildi", out, nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	raw, ok := c.Get("marketplace_user_id")
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token topilmadi", nil, nil)
		return
	}
	id, ok := raw.(uint)
	if !ok || id == 0 {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	out, err := h.service.GetMe(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "ok", out, nil)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	id, ok := getMarketplaceUserIDFromContext(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.UpdateMe(id, service.UpdateMeInput{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Gender:     req.Gender,
		RegionID:   req.RegionID,
		DistrictID: req.DistrictID,
		MFYID:      req.MFYID,
		BirthDate:  req.BirthDate,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "profil yangilandi", out, nil)
}

func (h *AuthHandler) GetMyAvatar(c *gin.Context) {
	id, ok := getMarketplaceUserIDFromContext(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	avatar, hasAvatar, err := h.service.GetMyAvatar(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "ok", gin.H{
		"has_avatar": hasAvatar,
		"avatar":     avatar,
	}, nil)
}

func (h *AuthHandler) UpdateMyAvatar(c *gin.Context) {
	id, ok := getMarketplaceUserIDFromContext(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	var req avatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "xato so'rov", nil, err.Error())
		return
	}
	out, err := h.service.UpdateMyAvatar(id, req.Avatar)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "avatar yangilandi", out, nil)
}

func (h *AuthHandler) DeleteMyAvatar(c *gin.Context) {
	id, ok := getMarketplaceUserIDFromContext(c)
	if !ok {
		response.JSON(c, http.StatusUnauthorized, "token yaroqsiz", nil, nil)
		return
	}
	out, err := h.service.DeleteMyAvatar(id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, "avatar o'chirildi", out, nil)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidPhone),
		errors.Is(err, service.ErrPurposeInvalid),
		errors.Is(err, service.ErrCodeInvalid),
		errors.Is(err, service.ErrInvalidGender),
		errors.Is(err, service.ErrFirstNameRequired),
		errors.Is(err, service.ErrLastNameRequired),
		errors.Is(err, service.ErrLocationInvalid),
		errors.Is(err, service.ErrLocationIDsRequired),
		errors.Is(err, service.ErrBirthDateInvalid),
		errors.Is(err, service.ErrAvatarBase64Invalid):
		response.JSON(c, http.StatusBadRequest, err.Error(), nil, nil)
	case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrCodeNotFound):
		response.JSON(c, http.StatusNotFound, err.Error(), nil, nil)
	case errors.Is(err, service.ErrUserAlreadyExists):
		response.JSON(c, http.StatusConflict, err.Error(), nil, nil)
	case errors.Is(err, service.ErrCodeExpired), errors.Is(err, service.ErrCodeNotVerified):
		response.JSON(c, http.StatusForbidden, err.Error(), nil, nil)
	default:
		response.JSON(c, http.StatusInternalServerError, strings.TrimSpace(err.Error()), nil, nil)
	}
}

func getMarketplaceUserIDFromContext(c *gin.Context) (uint, bool) {
	raw, ok := c.Get("marketplace_user_id")
	if !ok {
		return 0, false
	}
	id, ok := raw.(uint)
	if !ok || id == 0 {
		return 0, false
	}
	return id, true
}
