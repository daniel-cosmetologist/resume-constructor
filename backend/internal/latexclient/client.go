package latexclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client инкапсулирует HTTP-взаимодействие с latex-service.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
}

// NewClient создаёт новый клиент для общения с latex-service.
func NewClient(baseURL string, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.Default()
	}
	trimmed := strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL: trimmed,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: logger,
	}
}

// RenderPDF отправляет payload (данные резюме) в latex-service и возвращает PDF.
func (c *Client) RenderPDF(ctx context.Context, payload any) ([]byte, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("latex-service base URL is empty")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := c.baseURL + "/internal/v1/render"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to latex-service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limited := io.LimitReader(resp.Body, 4096)
		msg, _ := io.ReadAll(limited)
		c.logger.Printf("latex-service error: status=%d body=%s", resp.StatusCode, string(msg))
		return nil, fmt.Errorf("latex-service responded with status %d", resp.StatusCode)
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF response: %w", err)
	}

	if len(pdf) == 0 {
		return nil, fmt.Errorf("latex-service returned empty PDF")
	}

	return pdf, nil
}
