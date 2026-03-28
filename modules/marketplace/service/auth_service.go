package service

import (
	"context"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"time"

	"backend/internal/pkg/security"
	adminDomain "backend/modules/admin/domain"
	"backend/modules/eskiz"
	"backend/modules/marketplace/domain"
	"backend/modules/marketplace/repository"
)

var phoneRegex = regexp.MustCompile(`^\+998\d{9}$`)

var (
	ErrInvalidPhone        = errors.New("telefon raqami formati noto'g'ri")
	ErrUserNotFound        = errors.New("foydalanuvchi topilmadi")
	ErrUserAlreadyExists   = errors.New("bu telefon bo'yicha profil allaqachon mavjud")
	ErrInvalidGender       = errors.New("jins noto'g'ri")
	ErrFirstNameRequired   = errors.New("ism majburiy")
	ErrLastNameRequired    = errors.New("familiya majburiy")
	ErrLocationInvalid     = errors.New("viloyat, tuman va mfy mos emas")
	ErrLocationIDsRequired = errors.New("viloyat_id, tuman_id va mfy_id majburiy")
	ErrBirthDateInvalid    = errors.New("tug'ilgan sana noto'g'ri")
	ErrCodeNotFound        = errors.New("tasdiqlash kodi topilmadi")
	ErrCodeInvalid         = errors.New("tasdiqlash kodi noto'g'ri")
	ErrCodeExpired         = errors.New("tasdiqlash kodi muddati o'tgan")
	ErrCodeNotVerified     = errors.New("avval kodni tasdiqlang")
	ErrPurposeInvalid      = errors.New("purpose noto'g'ri")
	ErrAvatarBase64Invalid = errors.New("avatar base64 formati noto'g'ri")
)

type CheckPhoneResult struct {
	Exists  bool         `json:"exists"`
	Profile *UserProfile `json:"profile,omitempty"`
}

type UserProfile struct {
	ID         uint      `json:"id"`
	Phone      string    `json:"phone"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Gender     string    `json:"gender"`
	HasAvatar  bool      `json:"has_avatar"`
	Avatar     string    `json:"avatar,omitempty"`
	RegionID   uint      `json:"region_id"`
	DistrictID uint      `json:"district_id"`
	MFYID      uint      `json:"mfy_id"`
	BirthDate  string    `json:"birth_date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RegisterInput struct {
	Phone      string `json:"phone"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	BirthDate  string `json:"birth_date"` // YYYY-MM-DD
}

type VerifyCodeInput struct {
	Phone   string `json:"phone"`
	Code    string `json:"code"`
	Purpose string `json:"purpose"`
}

type SendCodeInput struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
}

type LoginInput struct {
	Phone string `json:"phone"`
}

type AuthResult struct {
	Token   string      `json:"token"`
	Profile UserProfile `json:"profile"`
}

type UpdateMeInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	RegionID   uint   `json:"region_id"`
	DistrictID uint   `json:"district_id"`
	MFYID      uint   `json:"mfy_id"`
	BirthDate  string `json:"birth_date"`
}

type EntryStartResult struct {
	Flow    string       `json:"flow"`
	Exists  bool         `json:"exists"`
	Sent    bool         `json:"sent"`
	Profile *UserProfile `json:"profile,omitempty"`
}

type EntryVerifyResult struct {
	Flow              string       `json:"flow"`
	Verified          bool         `json:"verified"`
	NeedsRegistration bool         `json:"needs_registration"`
	Token             string       `json:"token,omitempty"`
	Profile           *UserProfile `json:"profile,omitempty"`
}

type AuthService interface {
	EntryStart(phone string) (*EntryStartResult, error)
	EntryVerify(phone, code string) (*EntryVerifyResult, error)
	CheckPhone(phone string) (*CheckPhoneResult, error)
	SendCode(input SendCodeInput) error
	ResendCode(input SendCodeInput) error
	VerifyCode(input VerifyCodeInput) error
	LoginBySMS(input LoginInput) (*AuthResult, error)
	Register(input RegisterInput) (*AuthResult, error)
	GetMe(userID uint) (*UserProfile, error)
	UpdateMe(userID uint, input UpdateMeInput) (*UserProfile, error)
	GetMyAvatar(userID uint) (string, bool, error)
	UpdateMyAvatar(userID uint, avatarBase64 string) (*UserProfile, error)
	DeleteMyAvatar(userID uint) (*UserProfile, error)
}

type authService struct {
	repo           repository.AuthRepository
	eskizService   *eskiz.Service
	jwtSecret      string
	jwtExpireHours int
}

func NewAuthService(repo repository.AuthRepository, eskizService *eskiz.Service, jwtSecret string, jwtExpireHours int) AuthService {
	return &authService{repo: repo, eskizService: eskizService, jwtSecret: jwtSecret, jwtExpireHours: jwtExpireHours}
}

func (s *authService) EntryStart(phone string) (*EntryStartResult, error) {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err = s.SendCode(SendCodeInput{Phone: phone, Purpose: domain.CodePurposeLogin}); err != nil {
			return nil, err
		}
		return &EntryStartResult{
			Flow:    domain.CodePurposeLogin,
			Exists:  true,
			Sent:    true,
			Profile: toProfile(user),
		}, nil
	}
	if err = s.SendCode(SendCodeInput{Phone: phone, Purpose: domain.CodePurposeRegister}); err != nil {
		return nil, err
	}
	return &EntryStartResult{
		Flow:   domain.CodePurposeRegister,
		Exists: false,
		Sent:   true,
	}, nil
}

func (s *authService) EntryVerify(phone, code string) (*EntryVerifyResult, error) {
	phone = strings.TrimSpace(phone)
	code = strings.TrimSpace(code)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	if code == "" {
		return nil, ErrCodeInvalid
	}
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if err = s.VerifyCode(VerifyCodeInput{
			Phone:   phone,
			Code:    code,
			Purpose: domain.CodePurposeLogin,
		}); err != nil {
			return nil, err
		}
		auth, err := s.LoginBySMS(LoginInput{Phone: phone})
		if err != nil {
			return nil, err
		}
		profile := auth.Profile
		return &EntryVerifyResult{
			Flow:              domain.CodePurposeLogin,
			Verified:          true,
			NeedsRegistration: false,
			Token:             auth.Token,
			Profile:           &profile,
		}, nil
	}
	if err = s.VerifyCode(VerifyCodeInput{
		Phone:   phone,
		Code:    code,
		Purpose: domain.CodePurposeRegister,
	}); err != nil {
		return nil, err
	}
	return &EntryVerifyResult{
		Flow:              domain.CodePurposeRegister,
		Verified:          true,
		NeedsRegistration: true,
	}, nil
}

func (s *authService) CheckPhone(phone string) (*CheckPhoneResult, error) {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	row, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return &CheckPhoneResult{Exists: false}, nil
	}
	return &CheckPhoneResult{
		Exists:  true,
		Profile: toProfile(row),
	}, nil
}

func (s *authService) SendCode(input SendCodeInput) error {
	input.Phone = strings.TrimSpace(input.Phone)
	input.Purpose = strings.TrimSpace(input.Purpose)
	if !phoneRegex.MatchString(input.Phone) {
		return ErrInvalidPhone
	}
	if !validPurpose(input.Purpose) {
		return ErrPurposeInvalid
	}
	user, err := s.repo.GetUserByPhone(input.Phone)
	if err != nil {
		return err
	}
	if input.Purpose == domain.CodePurposeLogin && user == nil {
		return ErrUserNotFound
	}
	if input.Purpose == domain.CodePurposeRegister && user != nil {
		return ErrUserAlreadyExists
	}

	code := s.eskizService.GenerateCode()
	switch input.Purpose {
	case domain.CodePurposeRegister:
		_, err = s.eskizService.SendRegistrationCode(context.Background(), input.Phone, code)
	case domain.CodePurposeLogin:
		_, err = s.eskizService.SendLoginCode(context.Background(), input.Phone, code)
	}
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(5 * time.Minute)
	var userID *uint
	if user != nil {
		userID = &user.ID
	}
	row := &domain.VerificationCode{
		UserID:    userID,
		Phone:     input.Phone,
		Code:      code,
		Purpose:   input.Purpose,
		ExpiresAt: expiresAt,
	}
	return s.repo.CreateCode(row)
}

func (s *authService) ResendCode(input SendCodeInput) error {
	return s.SendCode(input)
}

func (s *authService) VerifyCode(input VerifyCodeInput) error {
	input.Phone = strings.TrimSpace(input.Phone)
	input.Code = strings.TrimSpace(input.Code)
	input.Purpose = strings.TrimSpace(input.Purpose)
	if !phoneRegex.MatchString(input.Phone) {
		return ErrInvalidPhone
	}
	if !validPurpose(input.Purpose) {
		return ErrPurposeInvalid
	}
	if input.Code == "" {
		return ErrCodeInvalid
	}
	row, err := s.repo.GetLatestActiveCode(input.Phone, input.Purpose)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrCodeNotFound
	}
	if time.Now().After(row.ExpiresAt) {
		return ErrCodeExpired
	}
	if row.Code != input.Code {
		row.Attempts++
		_ = s.repo.UpdateCode(row)
		return ErrCodeInvalid
	}
	now := time.Now()
	row.VerifiedAt = &now
	return s.repo.UpdateCode(row)
}

func (s *authService) LoginBySMS(input LoginInput) (*AuthResult, error) {
	phone := strings.TrimSpace(input.Phone)
	if !phoneRegex.MatchString(phone) {
		return nil, ErrInvalidPhone
	}
	user, err := s.repo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	code, err := s.repo.GetLatestActiveCode(phone, domain.CodePurposeLogin)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, ErrCodeNotFound
	}
	if time.Now().After(code.ExpiresAt) {
		return nil, ErrCodeExpired
	}
	if code.VerifiedAt == nil {
		return nil, ErrCodeNotVerified
	}
	now := time.Now()
	code.UsedAt = &now
	_ = s.repo.UpdateCode(code)
	token, err := security.GenerateMarketplaceToken(s.jwtSecret, user.ID, s.jwtExpireHours)
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, Profile: *toProfile(user)}, nil
}

func (s *authService) Register(input RegisterInput) (*AuthResult, error) {
	input.normalize()
	if !phoneRegex.MatchString(input.Phone) {
		return nil, ErrInvalidPhone
	}
	if input.FirstName == "" {
		return nil, ErrFirstNameRequired
	}
	if input.LastName == "" {
		return nil, ErrLastNameRequired
	}
	if input.Gender != domain.GenderMale && input.Gender != domain.GenderFemale {
		return nil, ErrInvalidGender
	}
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, ErrBirthDateInvalid
	}
	if err = s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	existing, err := s.repo.GetUserByPhone(input.Phone)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}
	code, err := s.repo.GetLatestActiveCode(input.Phone, domain.CodePurposeRegister)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, ErrCodeNotFound
	}
	if time.Now().After(code.ExpiresAt) {
		return nil, ErrCodeExpired
	}
	if code.VerifiedAt == nil {
		return nil, ErrCodeNotVerified
	}
	user := &domain.User{
		Phone:      input.Phone,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Gender:     input.Gender,
		RegionID:   input.RegionID,
		DistrictID: input.DistrictID,
		MFYID:      input.MFYID,
		BirthDate:  birthDate,
		Status:     adminDomain.StatusActive,
	}
	if err = s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	now := time.Now()
	code.UsedAt = &now
	_ = s.repo.UpdateCode(code)
	token, err := security.GenerateMarketplaceToken(s.jwtSecret, user.ID, s.jwtExpireHours)
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, Profile: *toProfile(user)}, nil
}

func (s *authService) GetMe(userID uint) (*UserProfile, error) {
	row, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	return toProfile(row), nil
}

func (s *authService) UpdateMe(userID uint, input UpdateMeInput) (*UserProfile, error) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Gender = strings.TrimSpace(input.Gender)
	input.BirthDate = strings.TrimSpace(input.BirthDate)

	if input.FirstName == "" {
		return nil, ErrFirstNameRequired
	}
	if input.LastName == "" {
		return nil, ErrLastNameRequired
	}
	if input.Gender != domain.GenderMale && input.Gender != domain.GenderFemale {
		return nil, ErrInvalidGender
	}
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, ErrBirthDateInvalid
	}
	if err = s.validateLocation(input.RegionID, input.DistrictID, input.MFYID); err != nil {
		return nil, err
	}
	row, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrUserNotFound
	}

	row.FirstName = input.FirstName
	row.LastName = input.LastName
	row.Gender = input.Gender
	row.RegionID = input.RegionID
	row.DistrictID = input.DistrictID
	row.MFYID = input.MFYID
	row.BirthDate = birthDate
	if err = s.repo.UpdateUser(row); err != nil {
		return nil, err
	}
	return toProfile(row), nil
}

func (s *authService) GetMyAvatar(userID uint) (string, bool, error) {
	row, err := s.repo.GetUserByID(userID)
	if err != nil {
		return "", false, err
	}
	if row == nil {
		return "", false, ErrUserNotFound
	}
	avatar := strings.TrimSpace(row.Avatar)
	return avatar, avatar != "", nil
}

func (s *authService) UpdateMyAvatar(userID uint, avatarBase64 string) (*UserProfile, error) {
	avatarBase64 = strings.TrimSpace(avatarBase64)
	if avatarBase64 == "" || !isValidBase64Image(avatarBase64) {
		return nil, ErrAvatarBase64Invalid
	}
	row, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	row.Avatar = avatarBase64
	if err = s.repo.UpdateUser(row); err != nil {
		return nil, err
	}
	return toProfile(row), nil
}

func (s *authService) DeleteMyAvatar(userID uint) (*UserProfile, error) {
	row, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrUserNotFound
	}
	row.Avatar = ""
	if err = s.repo.UpdateUser(row); err != nil {
		return nil, err
	}
	return toProfile(row), nil
}

func (s *authService) validateLocation(regionID, districtID, mfyID uint) error {
	if regionID == 0 || districtID == 0 || mfyID == 0 {
		return ErrLocationIDsRequired
	}
	district, err := s.repo.GetDistrictByID(districtID)
	if err != nil {
		return err
	}
	if district == nil || district.RegionID != regionID {
		return ErrLocationInvalid
	}
	mfy, err := s.repo.GetMFYByID(mfyID)
	if err != nil {
		return err
	}
	if mfy == nil || mfy.DistrictID != districtID {
		return ErrLocationInvalid
	}
	return nil
}

func validPurpose(p string) bool {
	return p == domain.CodePurposeLogin || p == domain.CodePurposeRegister
}

func toProfile(row *domain.User) *UserProfile {
	avatar := strings.TrimSpace(row.Avatar)
	return &UserProfile{
		ID:         row.ID,
		Phone:      row.Phone,
		FirstName:  row.FirstName,
		LastName:   row.LastName,
		Gender:     row.Gender,
		HasAvatar:  avatar != "",
		Avatar:     avatar,
		RegionID:   row.RegionID,
		DistrictID: row.DistrictID,
		MFYID:      row.MFYID,
		BirthDate:  row.BirthDate.Format("2006-01-02"),
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func (i *RegisterInput) normalize() {
	i.Phone = strings.TrimSpace(i.Phone)
	i.FirstName = strings.TrimSpace(i.FirstName)
	i.LastName = strings.TrimSpace(i.LastName)
	i.Gender = strings.TrimSpace(i.Gender)
	i.BirthDate = strings.TrimSpace(i.BirthDate)
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
