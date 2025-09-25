package ai

import (
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

func (q *Ai) buildPayload(message, systemPrompt string, survey json.RawMessage) ([]byte, error) {
	payload := sendPayload{
		SystemInstruction: &content{
			Parts: []part{{Text: systemPrompt}},
		},
		Contents: []content{
			{
				Role: aiRoleUser,
				Parts: []part{
					{Text: fmt.Sprintf(sendSurvey, string(survey))},
					{Text: message},
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

func BuildOutput[T any](reader io.Reader) (T, error) {
	var zero T

	body, err := io.ReadAll(reader)
	if err != nil {
		return zero, err
	}

	var resp respPayload
	if err = json.Unmarshal(body, &resp); err != nil {
		return zero, fmt.Errorf("decode ollama body: %w; body=%s", err, string(body))
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return zero, fmt.Errorf("empty candidates from gemini; body=%s", string(body))
	}

	raw := resp.Candidates[0].Content.Parts[0].Text
	var output T
	if err = json.Unmarshal([]byte(raw), &output); err != nil {
		logrus.Error("unmarshal model json: %w; raw=%s", err, raw)
		return zero, fmt.Errorf("unmarshal model json: %w; raw=%s", err, raw)
	}
	return output, nil
}
