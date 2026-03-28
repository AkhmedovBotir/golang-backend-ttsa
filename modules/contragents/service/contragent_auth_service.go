package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"backend/internal/pkg/security"
	adminDomain "backend/modules/admin/domain"
	"backend/modules/contragents/domain"
	"backend/modules/contragents/repository"
	"backend/modules/eskiz"
)

const codePurposePasswordSetup = "password_setup"

var phoneRegex = regexp.MustCompile(`^\+998\d{9}$`)

var (
	ErrInvalidPhone       = errors.New("telefon raqami formati noto'g'ri")
	ErrContragentNotFound = errors.New("kontragent topilmadi")
	ErrContragentInactive = errors.New("kontragent faol emas")
	ErrPasswordNotSet     = errors.New("parol hali o'rnatilmagan")
	ErrPasswordAlreadySet = errors.New("parol o'rnatilgan")
	ErrInvalidCredentials = errors.New("telefon yoki parol noto'g'ri")
	ErrCodeNotFound       = errors.New("tasdiqlash kodi topilmadi")
	ErrCodeExpired        = errors.New("tasdiqlash kodi muddati o'tgan")
	ErrCodeInvalid        = errors.New("tasdiqlash kodi noto'g'ri")
	ErrCodeNotVerified    = errors.New("avval kodni tasdiqlang")
	ErrPasswordTooShort   = errors.New("parol kamida 6 ta belgidan iborat bo'lishi kerak")
	ErrInvalidLogoBase64  = errors.New("logo base64 formati noto'g'ri")
)

type LoginResult struct {
	Token      string                  `json:"token"`
	Contragent *adminDomain.Contragent `json:"contragent"`
}

type ContragentAuthService interface {
	SendCode(phone string) error
	VerifyCode(phone, code string) error
	ResendCode(phone string) error
	SetPassword(phone, password string) (*LoginResult, error)
	Login(phone, password string) (*LoginResult, error)
	GetProfile(contragentID uint) (*adminDomain.Contragent, error)
	ChangePassword(contragentID uint, oldPassword, newPassword string) error
	UpdateLogo(contragentID uint, logoBase64 string) (*adminDomain.Contragent, error)
}

type contragentAuthService struct {
	repo           repository.ContragentAuthRepository
	eskizService   *eskiz.Service
	jwtSecret      string
	jwtExpireHours int
}

func NewContragentAuthService(
	repo repository.ContragentAuthRepository,
	eskizService *eskiz.Service,
	jwtSecret string,
	jwtExpireHours int,
) ContragentAuthService {
	return &contragentAuthService{
		repo:           repo,
		eskizService:   eskizService,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

func (s *contragentAuthService) SendCode(phone string) error {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}

	ctr, err := s.repo.GetContragentByPhone(phone)
	if err != nil {
		return err
	}
	if ctr == nil {
		return ErrContragentNotFound
	}
	if strings.ToLower(ctr.Status) != "active" {
		return ErrContragentInactive
	}
	if strings.TrimSpace(ctr.Password) != "" || !ctr.PasswordSetupAllowed {
		return ErrPasswordAlreadySet
	}

	code := s.eskizService.GenerateCode()
	message := fmt.Sprintf("%s - Kontragent hisobi uchun parol o'rnatish kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	if _, err = s.eskizService.SendSMS(context.Background(), phone, message); err != nil {
		return err
	}

	expires := time.Now().Add(5 * time.Minute)
	row := &domain.VerificationCode{
		ContragentID: ctr.ID,
		Phone:        phone,
		Code:         code,
		Purpose:      codePurposePasswordSetup,
		ExpiresAt:    expires,
	}
	return s.repo.CreateCode(row)
}

func (s *contragentAuthService) VerifyCode(phone, code string) error {
	phone = strings.TrimSpace(phone)
	code = strings.TrimSpace(code)
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	if code == "" {
		return ErrCodeInvalid
	}

	row, err := s.repo.GetLatestActiveCode(phone, codePurposePasswordSetup)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrCodeNotFound
	}
	if time.Now().After(row.ExpiresAt) {
		return ErrCodeExpired
	}
	if row.Code != code {
		row.Attempts++
		_ = s.repo.UpdateCode(row)
		return ErrCodeInvalid
	}

	now := time.Now()
	row.VerifiedAt = &now
	return s.repo.UpdateCode(row)
}

func (s *contragentAuthService) ResendCode(phone string) error {
	return s.SendCode(phone)
}

func (s *contragentAuthService) SetPassword(phone, password string) (*LoginResult, error) {
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	if len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	ctr, err := s.repo.GetContragentByPhone(phone)
	if err != nil {
		return nil, err
	}
	if ctr == nil {
		return nil, ErrContragentNotFound
	}
	if strings.TrimSpace(ctr.Password) != "" || !ctr.PasswordSetupAllowed {
		return nil, ErrPasswordAlreadySet
	}

	row, err := s.repo.GetLatestActiveCode(phone, codePurposePasswordSetup)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrCodeNotFound
	}
	if time.Now().After(row.ExpiresAt) {
		return nil, ErrCodeExpired
	}
	if row.VerifiedAt == nil {
		return nil, ErrCodeNotVerified
	}

	hashed, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}
	ctr.Password = hashed
	ctr.PasswordSetupAllowed = false
	if err = s.repo.UpdateContragent(ctr); err != nil {
		return nil, err
	}

	now := time.Now()
	row.UsedAt = &now
	if err = s.repo.UpdateCode(row); err != nil {
		return nil, err
	}

	token, err := security.GenerateContragentToken(s.jwtSecret, ctr.ID, s.jwtExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, Contragent: ctr}, nil
}

func (s *contragentAuthService) Login(phone, password string) (*LoginResult, error) {
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}

	ctr, err := s.repo.GetContragentByPhone(phone)
	if err != nil {
		return nil, err
	}
	if ctr == nil {
		return nil, ErrInvalidCredentials
	}
	if strings.ToLower(ctr.Status) != "active" {
		return nil, ErrContragentInactive
	}
	if strings.TrimSpace(ctr.Password) == "" {
		return nil, ErrPasswordNotSet
	}
	if !security.CheckPassword(ctr.Password, password) {
		return nil, ErrInvalidCredentials
	}

	token, err := security.GenerateContragentToken(s.jwtSecret, ctr.ID, s.jwtExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, Contragent: ctr}, nil
}

func (s *contragentAuthService) GetProfile(contragentID uint) (*adminDomain.Contragent, error) {
	row, err := s.repo.GetContragentByID(contragentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentNotFound
	}
	return row, nil
}

func (s *contragentAuthService) ChangePassword(contragentID uint, oldPassword, newPassword string) error {
	oldPassword = strings.TrimSpace(oldPassword)
	newPassword = strings.TrimSpace(newPassword)
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	row, err := s.repo.GetContragentByID(contragentID)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrContragentNotFound
	}
	if strings.TrimSpace(row.Password) == "" {
		return ErrPasswordNotSet
	}
	if !security.CheckPassword(row.Password, oldPassword) {
		return ErrInvalidCredentials
	}

	hashed, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	row.Password = hashed
	return s.repo.UpdateContragent(row)
}

func (s *contragentAuthService) UpdateLogo(contragentID uint, logoBase64 string) (*adminDomain.Contragent, error) {
	logoBase64 = strings.TrimSpace(logoBase64)
	if logoBase64 == "" {
		return nil, ErrInvalidLogoBase64
	}
	if !isValidBase64Image(logoBase64) {
		return nil, ErrInvalidLogoBase64
	}

	row, err := s.repo.GetContragentByID(contragentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrContragentNotFound
	}

	row.Logo = logoBase64
	if err = s.repo.UpdateContragent(row); err != nil {
		return nil, err
	}
	return row, nil
}

func isValidBase64Image(raw string) bool {
	payload := raw
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return false
		}
		payload = parts[1]
	}
	if payload == "" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(payload)
	return err == nil
}
