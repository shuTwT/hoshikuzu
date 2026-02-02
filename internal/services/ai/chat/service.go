package ai

import (
	setting_service "github.com/shuTwT/hoshikuzu/internal/services/system/setting"

	openai "github.com/sashabaranov/go-openai"
)

type AIService interface {
}

type AIServiceImpl struct {
	client         *openai.Client
	settingService setting_service.SettingService
}

func NewAIServiceImpl(settingService setting_service.SettingService) AIService {
	client := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL: "",
	})
	return &AIServiceImpl{
		client:         client,
		settingService: settingService,
	}
}
