package command

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
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
	gunslingerRepo      domain.GunslingerRepository
	userRepo            domain.UserRepository
	trans               *translator.Translator
	invited             map[uuid.UUID]map[int64]bool
	mu                  sync.Mutex
}

func NewJoinHandler(
	bot *service.Bot,
	announcer *service.GameAnnouncer,
	createGameUC *usecase.CreateGameUseCase,
	joinGameUC *usecase.JoinGameUseCase,
	calculateDeadlineUC *usecase.CalculateGameDeadlineUseCase,
	playGameUC *usecase.PlayGameUseCase,
	gunslingerRepo domain.GunslingerRepository,
	userRepo domain.UserRepository,
	trans *translator.Translator,
) Handler {
	return &JoinHandler{
		bot:                 bot,
		announcer:           announcer,
		createGameUC:        createGameUC,
		joinGameUC:          joinGameUC,
		calculateDeadlineUC: calculateDeadlineUC,
		playGameUC:          playGameUC,
		gunslingerRepo:      gunslingerRepo,
		userRepo:            userRepo,
		trans:               trans,
		invited:             make(map[uuid.UUID]map[int64]bool),
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
			message := h.trans.Get("player already in game", translator.Config{})

			if rand.Intn(2) == 0 {
				message += h.buildInviteSuffix(chatID, game)
			}

			h.bot.SendMessage(chatID, message)
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

func (h *JoinHandler) buildInviteSuffix(chatID int64, game *domain.Game) string {
	allPlayerIDs, err := h.gunslingerRepo.GetDistinctPlayerIDsInChat(chatID)
	if err != nil || len(allPlayerIDs) == 0 {
		return ""
	}

	currentPlayers := make(map[int64]bool, len(game.Gunslingers))
	for _, gs := range game.Gunslingers {
		currentPlayers[gs.PlayerID] = true
	}

	h.mu.Lock()
	alreadyInvited := h.invited[game.ID]
	if alreadyInvited == nil {
		alreadyInvited = make(map[int64]bool)
		h.invited[game.ID] = alreadyInvited
	}
	h.mu.Unlock()

	var candidates []int64
	for _, pid := range allPlayerIDs {
		if !currentPlayers[pid] && !alreadyInvited[pid] {
			candidates = append(candidates, pid)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	targetID := candidates[rand.Intn(len(candidates))]

	h.mu.Lock()
	alreadyInvited[targetID] = true
	h.mu.Unlock()

	user, err := h.userRepo.Get(targetID)
	if err != nil {
		return ""
	}

	return h.trans.Get("player already in game invite", translator.Config{
		Args: map[string]string{"%player": user.Mention()},
	})
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
