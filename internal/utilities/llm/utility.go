package llm

import "context"

type GeneratedMCQ struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

type GenerateQuestionsResponse struct {
	Questions []GeneratedMCQ `json:"questions"`
}

type GenerateQuestionsRequest struct {
	Title   string
	Content string
	Topic   string
	Tags    []string
}

type Utility interface {
	GenerateQuestions(ctx context.Context, request GenerateQuestionsRequest) (*GenerateQuestionsResponse, error)
}
