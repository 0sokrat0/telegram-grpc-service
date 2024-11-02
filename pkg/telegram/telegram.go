package telegram

import (
	"fmt"
	"github.com/0sokrat0/telegram-grpc-service/config"
	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func SendMessageToUser(userID int64, req proto_tg_service.SendMessageRequest) error {
	cfg := config.GetConfig()
	botToken := cfg.BotToken
	if botToken == "" {
		return fmt.Errorf("токен бота не указан в конфигурации")
	}

	if photoContent := req.GetPhotoContent(); photoContent != nil {
		// Подготовка параметров запроса для отправки фото
		params := url.Values{}
		params.Add("chat_id", strconv.FormatInt(userID, 10))
		params.Add("photo", photoContent.Url)
		if photoContent.Caption != "" {
			params.Add("caption", photoContent.Caption)
		}
		if photoContent.ParseMode != "" {
			params.Add("parse_mode", photoContent.ParseMode)
		}

		// Отправка запроса к Telegram API
		resp, err := http.PostForm("https://api.telegram.org/bot"+botToken+"/sendPhoto", params)
		if err != nil {
			return fmt.Errorf("ошибка при отправке запроса к Telegram API: %v", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("не удалось прочитать ответ от Telegram API: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("не удалось отправить фото, статус: %s, ответ: %s", resp.Status, string(bodyBytes))
		}

	} else if textContent := req.GetTextContent(); textContent != nil {
		// Обработка текстового сообщения...
		// (оставьте существующий код для обработки текстовых сообщений)
	} else {
		return fmt.Errorf("не указан контент для отправки")
	}

	return nil
}
