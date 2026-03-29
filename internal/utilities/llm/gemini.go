package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpclient"
)

type geminiUtility struct {
	httpClient httpclient.Utility
	config     config.Utility
}

// JSON structures for Gemini API
type generateContentRequest struct {
	Contents []content `json:"contents"`
}

type generateContentResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content content `json:"content"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

func (m *geminiUtility) GenerateQuestions(ctx context.Context, request GenerateQuestionsRequest) (*GenerateQuestionsResponse, error) {
	cfg := m.config.Config()

	if strings.TrimSpace(cfg.LLMBaseURL) == "" {
		return nil, errors.New("llm base url is required")
	}
	if strings.TrimSpace(cfg.LLMAPIKey) == "" {
		return nil, errors.New("llm api key is required")
	}

	model := strings.TrimSpace(cfg.LLMModel)
	if model == "" {
		model = "gemini-2.5-flash"
	}

	var response generateContentResponse
	_, err := m.httpClient.DoJSON(ctx, httpclient.Request{
		Method:  http.MethodPost,
		URL:     strings.TrimRight(cfg.LLMBaseURL, "/") + "/v1beta/models/" + model + ":generateContent",
		Timeout: 30 * time.Second,
		Headers: map[string]string{
			"x-goog-api-key": cfg.LLMAPIKey,
			"Content-Type":   "application/json",
		},
		Body: generateContentRequest{
			Contents: []content{
				{
					Parts: []part{
						{
							Text: buildQuestionsPrompt(request),
						},
					},
				},
			},
		},
	}, &response)
	if err != nil {
		return nil, err
	}

	return extractQuestions(response)
}

func buildQuestionsPrompt(request GenerateQuestionsRequest) string {
	return fmt.Sprintf(
		`Generate 5 multiple-choice study questions from the note below.

		Return ONLY valid JSON.
		Do not include markdown.
		Do not include code fences.
		Do not include explanation text.

		Use exactly this JSON structure:
		{
		"questions": [
			{
			"question": "string",
			"options": ["string", "string", "string", "string"],
			"answer": "string"
			}
		]
		}

		Rules:
		- Generate exactly 5 questions
		- Each question must have exactly 4 options
		- The answer must exactly match one of the options
		- Keep questions concise and clear
		- Base the questions only on the note content

		Title: %s
		Topic: %s
		Tags: %s
		Content:
		%s`,
		request.Title,
		request.Topic,
		strings.Join(request.Tags, ", "),
		request.Content,
	)
}

// extractQuestions parses the Gemini API response and validates the generated questions.
func extractQuestions(response generateContentResponse) (*GenerateQuestionsResponse, error) {
	// Basic validation to ensure response contains expected content
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("llm response did not contain any content")
	}

	text := strings.TrimSpace(response.Candidates[0].Content.Parts[0].Text)
	if text == "" {
		return nil, errors.New("llm response was empty")
	}

	// Unmarshal the cleaned text into our expected structure (question, options, answer)
	var result GenerateQuestionsResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal llm json response: %w", err)
	}

	if len(result.Questions) == 0 {
		return nil, errors.New("llm response did not contain any questions")
	}

	// Validate each question to ensure it meets our criteria
	for i, question := range result.Questions {
		if strings.TrimSpace(question.Question) == "" {
			return nil, fmt.Errorf("question %d is empty", i)
		}
		if len(question.Options) != 4 {
			return nil, fmt.Errorf("question %d does not contain exactly 4 options", i)
		}
		if strings.TrimSpace(question.Answer) == "" {
			return nil, fmt.Errorf("question %d answer is empty", i)
		}

		answerMatchesOption := false
		for _, option := range question.Options {
			if strings.TrimSpace(option) == strings.TrimSpace(question.Answer) {
				answerMatchesOption = true
				break
			}
		}
		if !answerMatchesOption {
			return nil, fmt.Errorf("question %d answer does not match any option", i)
		}
	}

	return &result, nil
}
