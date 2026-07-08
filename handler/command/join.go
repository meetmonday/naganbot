package command

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"github.com/taranovegor/naganbot/usecase"
)

type JoinHandler struct {
	Handler
	bot                 *service.Bot
	announcer           *service.GameAnnouncer
	createGameUC        *usecase.CreateGameUseCase
	joinGameUC          *usecase.JoinGameUseCase
	calculateDeadlineUC *usecase.CalculateGameDeadlineUseCase
	playGameUC          *usecase.PlayGameUseCase
	trans               *translator.Translator
}

func NewJoinHandler(
	bot *service.Bot,
	announcer *service.GameAnnouncer,
	createGameUC *usecase.CreateGameUseCase,
	joinGameUC *usecase.JoinGameUseCase,
	calculateDeadlineUC *usecase.CalculateGameDeadlineUseCase,
	playGameUC *usecase.PlayGameUseCase,
	trans *translator.Translator,
) Handler {
	return &JoinHandler{
		bot:                 bot,
		announcer:           announcer,
		createGameUC:        createGameUC,
		joinGameUC:          joinGameUC,
		calculateDeadlineUC: calculateDeadlineUC,
		playGameUC:          playGameUC,
		trans:               trans,
	}
}

func (h *JoinHandler) Name() string {
	return "join"
}

func (h *JoinHandler) Execute(msg *tgbotapi.Message) {
	chatID, userID := msg.Chat.ID, msg.From.ID
	game, err := h.createGameUC.Execute(chatID, userID)
	if err != nil {
		return
	}

	_, err = h.joinGameUC.Execute(game.ID, msg.From.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrPlayerAlreadyInGame) {
			h.bot.SendMessage(chatID, h.trans.Get("player already in game", translator.Config{}))
		}
		return
	}

	game, count, err := h.calculateDeadlineUC.Execute(game.ID, time.Now())
	if err != nil {
		fmt.Println(err)
		h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if game.Mode == domain.GameModeDynamic {
		h.handleDynamicJoin(chatID, game, count)
		return
	}

	h.handleClassicJoin(chatID, userID, game)
}

func (h *JoinHandler) handleClassicJoin(chatID int64, userID int64, game *domain.Game) {
	hitReport, err := h.playGameUC.Execute(context.TODO(), game.ID)
	if err != nil {
		fmt.Println(err)
		if game.Owner.ID != userID && errors.Is(err, usecase.ErrNotEnoughPlayers) {
			h.bot.SendMessage(chatID, h.trans.Get("joining the game", translator.Config{}))
		} else if !errors.Is(err, usecase.ErrNotEnoughPlayers) {
			h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
		}
		return
	}

	h.announcer.AnnounceGameResult(chatID, hitReport)
}

func (h *JoinHandler) handleDynamicJoin(chatID int64, game *domain.Game, count int) {
	if game.ShouldStartNow(count) {
		hitReport, err := h.playGameUC.Execute(context.TODO(), game.ID)
		if err != nil {
			fmt.Println(err)
			if !errors.Is(err, usecase.ErrNotEnoughPlayers) {
				h.bot.SendMessage(chatID, h.trans.Get("something went wrong", translator.Config{}))
			}
			return
		}

		h.announcer.AnnounceGameResult(chatID, hitReport)
		return
	}

	var deadline string
	if count >= domain.DynamicMinPlayers {
		deadline = h.formatDeadline(game.StartDeadline.Time)
	}

	h.bot.SendMessage(chatID, h.joinMessage(chatID, count, domain.DynamicMaxPlayers, domain.DynamicMinPlayers, deadline))
}

func (h *JoinHandler) joinMessage(chatID int64, count int, max int, min int, deadline string) string {
	base := h.trans.Get("joining the game", translator.Config{})

	var details string
	if deadline == "" {
		details = h.trans.Get("joining the game details waiting", translator.Config{
			Args: map[string]string{
				"%count": strconv.Itoa(count),
				"%max":   strconv.Itoa(max),
				"%min":   strconv.Itoa(min),
			},
		})
	} else {
		details = h.trans.Get("joining the game details deadline", translator.Config{
			Args: map[string]string{
				"%count":    strconv.Itoa(count),
				"%max":      strconv.Itoa(max),
				"%deadline": deadline,
			},
		})
	}

	return base + details
}

func (h *JoinHandler) formatDeadline(deadline time.Time) string {
	midnight := domain.NextMidnight(time.Now())
	if deadline.Equal(midnight) || (deadline.After(midnight.Add(-time.Minute)) && deadline.Before(midnight.Add(time.Minute))) {
		return h.trans.Get("game starts at midnight", translator.Config{})
	}

	minutes := int(time.Until(deadline).Minutes())
	if minutes < 2 {
		return h.trans.Get("game starts in less than a minute", translator.Config{})
	}

	return h.trans.Get("game starts in minutes", translator.Config{Count: minutes})
}
