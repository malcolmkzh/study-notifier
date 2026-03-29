package llm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/malcolmkzh/study-notifier/internal/utilities/config"
	"github.com/malcolmkzh/study-notifier/internal/utilities/httpclient"
)

func NewLLMUtility(configUtility config.Utility, httpClientUtility httpclient.Utility) (Utility, error) {
	if configUtility == nil {
		return nil, errors.New("config utility is required")
	}
	if httpClientUtility == nil {
		return nil, errors.New("http client utility is required")
	}

	provider := strings.TrimSpace(configUtility.Config().LLMProvider)
	if provider == "" {
		provider = "gemini"
	}

	switch strings.ToLower(provider) {
	case "gemini":
		return &geminiUtility{
			httpClient: httpClientUtility,
			config:     configUtility,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", provider)
	}
}
