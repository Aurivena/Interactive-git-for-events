package ai

import (
	"arch/internal/domain/entity"
	"bytes"
	"context"
	"errors"
	"io"

	"fmt"

	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	methodPost   = "POST"
	DefaultCount = 3
	url          = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"
)

type Ai struct {
	ai     entity.AiConfig
	client *http.Client
}

func New(cfg entity.AiConfig, httpClient *http.Client) *Ai {
	return &Ai{
		ai:     cfg,
		client: httpClient,
	}
}

func (q *Ai) Send(message string) ([]entity.RequestPayload, error) {
	payload, err := q.buildPayload(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(methodPost, fmt.Sprintf(url, q.ai.Model, q.ai.ApiKey), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "arch-ai-client/1.0")

	resp, err := q.client.Do(req)
	if err != nil {
		logrus.Error("error at send request from AI", err)
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("close body error: %v\n", err)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("499")
		}
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := string(b)
		logrus.Warn(resp.StatusCode)
		logrus.Error("error at send request from AI", msg)
		if len(msg) > 2_000 {
			msg = msg[:2_000] + "..."
		}
		return nil, fmt.Errorf("AI http %d: %s", resp.StatusCode, msg)
	}

	output, err := buildOutput(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return output, nil
}
