package eskiz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL     = "https://notify.eskiz.uz/api"
	defaultSender      = "4546"
	tokenRefreshPeriod = 29 * 24 * time.Hour
)

var phoneCleaner = regexp.MustCompile(`[+\s\-()]`)

var (
	ErrEskizConfigInvalid = errors.New("eskiz sozlamalari to'liq emas")
	ErrEskizAuthFailed    = errors.New("SMS xizmatiga ulanishda xatolik yuz berdi")
	ErrEskizSendFailed    = errors.New("SMS yuborishda xatolik yuz berdi")
)

type Service struct {
	email    string
	password string
	baseURL  string
	sender   string
	client   *http.Client

	mu             sync.Mutex
	token          string
	tokenExpiresAt time.Time
}

type SendResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Status    string `json:"status"`
}

func NewFromEnv() (*Service, error) {
	email := strings.TrimSpace(os.Getenv("ESKIZ_EMAIL"))
	password := strings.TrimSpace(os.Getenv("ESKIZ_PASSWORD"))
	baseURL := strings.TrimSpace(os.Getenv("ESKIZ_BASE_URL"))
	sender := strings.TrimSpace(os.Getenv("ESKIZ_SENDER"))
	timeoutRaw := strings.TrimSpace(os.Getenv("ESKIZ_TIMEOUT_SECONDS"))

	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if sender == "" {
		sender = defaultSender
	}

	timeout := 15 * time.Second
	if timeoutRaw != "" {
		if n, err := time.ParseDuration(timeoutRaw + "s"); err == nil && n > 0 {
			timeout = n
		}
	}

	return &Service{
		email:    email,
		password: password,
		baseURL:  strings.TrimRight(baseURL, "/"),
		sender:   sender,
		client:   &http.Client{Timeout: timeout},
	}, nil
}

func New(email, password string) *Service {
	return &Service{
		email:    strings.TrimSpace(email),
		password: strings.TrimSpace(password),
		baseURL:  defaultBaseURL,
		sender:   defaultSender,
		client:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *Service) GetToken(ctx context.Context) (string, error) {
	if strings.TrimSpace(s.email) == "" || strings.TrimSpace(s.password) == "" {
		return "", ErrEskizConfigInvalid
	}

	s.mu.Lock()
	if s.token != "" && !s.tokenExpiresAt.IsZero() && time.Now().Before(s.tokenExpiresAt) {
		token := s.token
		s.mu.Unlock()
		return token, nil
	}
	s.mu.Unlock()

	reqBody := map[string]string{
		"email":    s.email,
		"password": s.password,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", ErrEskizAuthFailed
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", ErrEskizAuthFailed
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(raw, &parsed)

	if resp.StatusCode >= http.StatusBadRequest || parsed.Data.Token == "" {
		return "", fmt.Errorf("%w: %s", ErrEskizAuthFailed, strings.TrimSpace(string(raw)))
	}

	s.mu.Lock()
	s.token = parsed.Data.Token
	s.tokenExpiresAt = time.Now().Add(tokenRefreshPeriod)
	token := s.token
	s.mu.Unlock()

	return token, nil
}

func (s *Service) SendSMS(ctx context.Context, phone, message string) (*SendResult, error) {
	return s.sendSMS(ctx, phone, message, true)
}

func (s *Service) sendSMS(ctx context.Context, phone, message string, allowRetry bool) (*SendResult, error) {
	token, err := s.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	finalPhone := normalizePhone(phone)
	reqBody := map[string]string{
		"mobile_phone": finalPhone,
		"message":      message,
		"from":         s.sender,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/message/sms/send", bytes.NewReader(body))
	if err != nil {
		return nil, ErrEskizSendFailed
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "waiting") {
			return &SendResult{Success: true, Status: "waiting"}, nil
		}
		return nil, ErrEskizSendFailed
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed struct {
		Status  string      `json:"status"`
		ID      interface{} `json:"id"`
		Message string      `json:"message"`
	}
	_ = json.Unmarshal(raw, &parsed)

	if resp.StatusCode == http.StatusUnauthorized && allowRetry {
		s.mu.Lock()
		s.token = ""
		s.tokenExpiresAt = time.Time{}
		s.mu.Unlock()
		return s.sendSMS(ctx, phone, message, false)
	}

	if strings.Contains(strings.ToLower(parsed.Message), "waiting") {
		return &SendResult{
			Success:   true,
			MessageID: toString(parsed.ID),
			Status:    "waiting",
		}, nil
	}

	if parsed.Status == "success" || parsed.Status == "waiting" || parsed.ID != nil {
		status := parsed.Status
		if status == "" {
			status = "sent"
		}
		return &SendResult{
			Success:   true,
			MessageID: toString(parsed.ID),
			Status:    status,
		}, nil
	}

	if parsed.Message != "" {
		return nil, fmt.Errorf("%w: %s", ErrEskizSendFailed, parsed.Message)
	}
	return &SendResult{
		Success:   true,
		MessageID: toString(parsed.ID),
		Status:    "sent",
	}, nil
}

func (s *Service) GenerateCode() string {
	n := 10000 + rand.Intn(90000)
	return fmt.Sprintf("%d", n)
}

func (s *Service) SendRegistrationCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Marketplace ilovasidan ro'yxatdan o'tish uchun tasdiqlash kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendLoginCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Marketplace ilovasiga kirish uchun tasdiqlash kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendForgotPasswordCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Marketplace ilovasida parol tiklash uchun tasdiqlash kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendContragentPasswordSetupCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Kontragent hisobi uchun parol o'rnatish kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendDeliveryProviderPasswordSetupCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Yetkazib beruvchi hisobi uchun parol o'rnatish kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendPunktPasswordSetupCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Punkt hisobi uchun parol o'rnatish kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendAgentPasswordSetupCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Agent hisobi uchun parol o'rnatish kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func (s *Service) SendDeviceVerificationCode(ctx context.Context, phone, code string) (*SendResult, error) {
	message := fmt.Sprintf("%s - Yangi qurilmani tasdiqlash uchun tasdiqlash kodi. Kod 5 daqiqa amal qiladi. Talab va Taklif Agency.", code)
	return s.SendSMS(ctx, phone, message)
}

func normalizePhone(phone string) string {
	formatted := phoneCleaner.ReplaceAllString(phone, "")
	if strings.HasPrefix(formatted, "998") {
		return formatted
	}
	return "998" + formatted
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
