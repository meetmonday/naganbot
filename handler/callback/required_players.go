package callback

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

var revolverOptions = []int{4, 5, 6, 7, 8, 9}

func RevolverKeyboard(mode domain.GameMode, selected int, trans *translator.Translator) service.InlineKeyboard {
	var keyboard service.InlineKeyboard
	for _, n := range revolverOptions {
		s := strconv.Itoa(n)
		arg := RequiredPlayers.SetArgs(s).ToString()
		txt := trans.Get(fmt.Sprintf("%s shot revolver", s), translator.Config{})
		if mode == domain.GameModeClassic && n == selected {
			txt = fmt.Sprintf("🔫 %s", txt)
		}
		keyboard = append(keyboard, []service.KeyboardButton{{Data: arg, Label: txt}})
	}

	dynamicArg := RequiredPlayers.SetArgs("dynamic").ToString()
	dynamicTxt := trans.Get("dynamic shot revolver", translator.Config{})
	if mode == domain.GameModeDynamic {
		dynamicTxt = fmt.Sprintf("🔫 %s", dynamicTxt)
	}
	keyboard = append(keyboard, []service.KeyboardButton{{Data: dynamicArg, Label: dynamicTxt}})

	return keyboard
}

type requiredPlayers struct {
	chatRepo domain.ChatRepository
	bot      *service.Bot
	trans    *translator.Translator
}

func NewRequiredPlayers(
	chatRepo domain.ChatRepository,
	bot *service.Bot,
	trans *translator.Translator,
) Handler {
	return &requiredPlayers{
		chatRepo: chatRepo,
		bot:      bot,
		trans:    trans,
	}
}

func (h *requiredPlayers) Pattern() Pattern {
	return RequiredPlayers
}

func (h *requiredPlayers) Execute(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID

	allowedUsername := config.GetEnv(config.ForceUsername)
	if allowedUsername != "" && query.From.UserName != allowedUsername {
		return
	}

	isAdmin, err := h.bot.IsAdmin(chatID, query.From.ID)
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if !isAdmin {
		notification := h.trans.Get("settings can be changed only by admins", translator.Config{})
		h.bot.AnswerCallback(query.ID, notification)
		return
	}

	arg := RequiredPlayers.GetArg(query.Data, 1)

	chat, err := h.chatRepo.Get(chatID)
	if err != nil {
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	var notification string
	if arg == "dynamic" {
		chat.Settings.Mode = domain.GameModeDynamic
		notification = fmt.Sprintf(
			"%s\n%s",
			h.trans.Get("dynamic mode enabled", translator.Config{}),
			h.trans.Get("settings will be applied for next games", translator.Config{}),
		)
	} else {
		players, err := strconv.Atoi(arg)
		if err != nil {
			h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
			return
		}

		chat.Settings.Mode = domain.GameModeClassic
		chat.Settings.RequiredPlayers = players
		notification = fmt.Sprintf(
			"%s\n%s",
			h.trans.Get("revolver has been replaced", translator.Config{Count: players}),
			h.trans.Get("settings will be applied for next games", translator.Config{}),
		)
	}

	h.chatRepo.Update(&chat)
	h.bot.AnswerCallback(query.ID, notification)

	keyboard := RevolverKeyboard(chat.Settings.Mode, chat.Settings.RequiredPlayers, h.trans)
	h.bot.EditMessageReplyMarkup(chatID, query.Message.MessageID, keyboard)
}
