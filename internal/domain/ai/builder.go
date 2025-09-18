package ai

import (
	"arch/internal/domain/entity"
	"encoding/json"
	"fmt"

	"io"

	"github.com/sirupsen/logrus"
)

const (
	aiRoleUser     = "user"
	maxOutputToken = 3000
)

type sendPayload struct {
	SystemInstruction *content          `json:"system_instruction,omitempty"`
	Contents          []content         `json:"contents"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	ResponseMimeType string `json:"response_mime_type,omitempty"`
	MaxOutputTokens  int32  `json:"maxOutputTokens,omitempty"`
}

type respPayload struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (q *Ai) buildPayload(message string) ([]byte, error) {
	payload := sendPayload{
		SystemInstruction: &content{
			Parts: []part{{Text: sendPrompt}},
		},
		Contents: []content{
			{
				Role: aiRoleUser,
				Parts: []part{
					{
						Text: message,
					},
				},
			},
		},
		GenerationConfig: &generationConfig{
			ResponseMimeType: "application/json",
			MaxOutputTokens:  maxOutputToken,
		},
	}

	return json.Marshal(payload)
}

func buildOutput(reader io.Reader) ([]entity.RequestPayload, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var resp respPayload
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode ollama body: %w; body=%s", err, string(body))
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty candidates from gemini; body=%s", string(body))
	}

	raw := resp.Candidates[0].Content.Parts[0].Text
	var output []entity.RequestPayload
	if err = json.Unmarshal([]byte(raw), &output); err != nil {
		logrus.Error("unmarshal model json: %w; raw=%s", err, raw)
		return nil, fmt.Errorf("unmarshal model json: %w; raw=%s", err, raw)
	}
	return output, nil
}
