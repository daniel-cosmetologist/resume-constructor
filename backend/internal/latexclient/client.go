package latexclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"resume_backend/internal/resume"
)

// Client реализует вызов LaTeX-сервиса по HTTP.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
}

func NewClient(baseURL string, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.Default()
	}
	if baseURL == "" {
		baseURL = "http://latex-service:8081"
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// RenderResume отправляет JSON с резюме в LaTeX-сервис и возвращает PDF.
func (c *Client) RenderResume(ctx context.Context, r resume.Resume) ([]byte, error) {
	url := fmt.Sprintf("%s/internal/v1/render", c.baseURL)

	payload, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshal resume: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call latex-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		c.logger.Printf("latex-service returned %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("latex-service returned status %d", resp.StatusCode)
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read pdf body: %w", err)
	}

	return pdf, nil
}
