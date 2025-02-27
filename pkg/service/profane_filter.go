package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ProfaneFilterService проверяет комментарии на запрещенные слова
type ProfaneFilterService struct {
	apiURL string
	client *http.Client
}

func NewProfaneFilterService(apiURL string) *ProfaneFilterService {
	return &ProfaneFilterService{
		apiURL: apiURL,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// ErrProfaneText используется, если текст не прошел валидацию
var ErrProfaneText = errors.New("profane text detected")

// Check проверяет текст на запрещенные слова
func (p *ProfaneFilterService) Check(text string, traceID string) error {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w", err)
	}

	filterURL, err := url.JoinPath(p.apiURL, "v1/filter/")
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequest("POST", filterURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", traceID) // Пробрасываем trace_id

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Обрабатываем статус-код 403 как ошибку валидации
	if resp.StatusCode == http.StatusForbidden {
		var errResponse struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("failed to parse 403 response: %w", err)
		}
		return fmt.Errorf("%w: %s", ErrProfaneText, errResponse.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	// Проверяем успешный ответ
	var successResponse struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		return fmt.Errorf("failed to parse success response: %w", err)
	}
	if successResponse.Text != text {
		return fmt.Errorf("text validation mismatch: received unexpected response")
	}

	return nil
}
