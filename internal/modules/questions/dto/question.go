package dto

import "github.com/malcolmkzh/study-notifier/internal/utilities/llm"

type GenerateQuestionsResponse struct {
	NoteID    uint               `json:"note_id"`
	Questions []llm.GeneratedMCQ `json:"questions"`
	Message   string             `json:"message"`
}
