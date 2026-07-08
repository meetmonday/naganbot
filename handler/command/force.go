package command

import (
	"context"
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type ForceHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	game       domain.GameRepository
	playGameUC *usecase.PlayGameUseCase
}

func NewForceHandler(
	bot *service.Bot,
	trans *translator.Translator,
	game domain.GameRepository,
	playGameUC *usecase.PlayGameUseCase,
) Handler {
	return &ForceHandler{
		bot:        bot,
		trans:      trans,
		game:       game,
		playGameUC: playGameUC,
	}
}

func (hdlr ForceHandler) Name() string {
	return "force"
}

func (hdlr ForceHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	allowedUsername := config.GetEnv(config.ForceUsername)
	if allowedUsername != "" && msg.From.UserName != allowedUsername {
		return
	}

	game, err := hdlr.game.GetActiveForChat(chatID)
	if err != nil {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("active game not found", translator.Config{}))
		return
	}

	hitReport, err := hdlr.playGameUC.ForceExecute(context.TODO(), game.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrGameAlreadyPlayed) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("active game not found", translator.Config{}))
		} else if errors.Is(err, usecase.ErrNotEnoughPlayers) {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("not enough players", translator.Config{}))
		} else {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
		}
		return
	}

	for _, message := range hdlr.trans.GetMany("play the game", translator.Config{}) {
		hdlr.bot.SendMessage(chatID, message)
		time.Sleep(time.Second)
	}

	isAtomic := hitReport.BulletType == service.BulletAtomicType

	if isAtomic {
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("killed by atomic bullet", translator.Config{}))
	}

	for _, victim := range hitReport.Victims {
		if !isAtomic {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
		}

		err = hdlr.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("player is not kicked", translator.Config{}))
		}
	}
}
