package service

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const parseMode = tgbotapi.ModeHTML

type KeyboardButton struct {
	Data  string
	Label string
}

type InlineKeyboard [][]KeyboardButton

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(
	api *tgbotapi.BotAPI,
) *Bot {
	return &Bot{
		api: api,
	}
}

func (bot Bot) SendMessage(chatID int64, text string) error {
	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: chatID},
		ParseMode:             parseMode,
		Text:                  text,
		DisableWebPagePreview: true,
	})

	return err
}

func (bot Bot) DeleteMessage(chatID int64, messageID int) error {
	_, err := bot.api.Request(tgbotapi.NewDeleteMessage(chatID, messageID))

	return err
}

func (bot Bot) Ban(chatID int64, userID int64, untilDate int64) error {
	_, err := bot.api.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: untilDate,
	})

	return err
}

func (bot Bot) Kick(chatID int64, userID int64) error {
	return bot.Ban(chatID, userID, time.Now().Add(time.Minute).Unix())
}

func (bot Bot) buildKeyboardMarkup(keyboard InlineKeyboard) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, row := range keyboard {
		var cols []tgbotapi.InlineKeyboardButton
		for _, btn := range row {
			cols = append(cols, tgbotapi.NewInlineKeyboardButtonData(btn.Label, btn.Data))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(cols...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (bot Bot) SendInlineKeyboard(chatID int64, text string, keyboard InlineKeyboard) error {
	_, err := bot.api.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: bot.buildKeyboardMarkup(keyboard),
		},
		ParseMode: parseMode,
		Text:      text,
	})

	return err
}

func (bot Bot) EditMessage(chatID int64, messageID int, text string, keyboard InlineKeyboard) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = parseMode
	markup := bot.buildKeyboardMarkup(keyboard)
	edit.ReplyMarkup = &markup
	_, err := bot.api.Request(edit)
	return err
}

func (bot Bot) EditMessageReplyMarkup(chatID int64, messageID int, keyboard InlineKeyboard) error {
	markup := bot.buildKeyboardMarkup(keyboard)
	_, err := bot.api.Request(tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, markup))

	return err
}

func (bot Bot) AnswerCallback(callbackQueryID string, text string) error {
	_, err := bot.api.Request(tgbotapi.NewCallback(callbackQueryID, text))

	return err
}

func (bot Bot) IsAdmin(chatID int64, userID int64) (bool, error) {
	resp, err := bot.api.Request(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	})
	if err != nil {
		return false, err
	}

	var records []struct {
		User struct {
			ID int64 `json:"id"`
		} `json:"user"`
	}

	err = json.Unmarshal(resp.Result, &records)
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record.User.ID == userID {
			return true, nil
		}
	}

	return false, nil
}
