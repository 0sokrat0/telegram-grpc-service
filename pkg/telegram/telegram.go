package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"telegram-grpc-service/config"
	proto_tg_service "telegram-grpc-service/gen/go/proto"
)

// SendMessageToUser отправляет сообщение пользователю через Telegram Bot API
func SendMessageToUser(userID int64, req proto_tg_service.SendMessageRequest) error {
	cfg := config.GetConfig()
	botToken := cfg.BotToken
	if botToken == "" {
		return fmt.Errorf("токен бота не указан в конфигурации")
	}

	// Определяем тип контента и подготавливаем запрос
	if textContent := req.GetTextContent(); textContent != nil {
		// Обработка TextContent
		params := url.Values{}
		params.Add("chat_id", strconv.FormatInt(userID, 10))
		params.Add("text", textContent.Text)
		if textContent.ParseMode != "" {
			params.Add("parse_mode", textContent.ParseMode)
		}
		if textContent.DisableWebPagePreview {
			params.Add("disable_web_page_preview", "true")
		}

		// Отправка запроса к Telegram API
		resp, err := http.PostForm("https://api.telegram.org/bot"+botToken+"/sendMessage", params)
		if err != nil {
			return fmt.Errorf("ошибка при отправке запроса к Telegram API: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("не удалось отправить сообщение, статус: %s", resp.Status)
		}

	} else if photoContent := req.GetPhotoContent(); photoContent != nil {
		// Обработка PhotoContent
		params := url.Values{}
		params.Add("chat_id", strconv.FormatInt(userID, 10))
		params.Add("photo", photoContent.Url)
		if photoContent.Caption != "" {
			params.Add("caption", photoContent.Caption)
		}

		// Отправка запроса к Telegram API
		resp, err := http.PostForm("https://api.telegram.org/bot"+botToken+"/sendPhoto", params)
		if err != nil {
			return fmt.Errorf("ошибка при отправке запроса к Telegram API: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("не удалось отправить фото, статус: %s", resp.Status)
		}

	} else {
		return fmt.Errorf("не указан контент для отправки")
	}

	return nil
}
