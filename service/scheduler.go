package service

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/config"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/translator"
)

const (
	SchedulerTickInterval   = time.Minute
	SchedulerStartInterval  = 2 * time.Second
	SchedulerMidnightHour   = 0
	SchedulerMidnightMinute = 0
)

type GameStarter interface {
	Execute(ctx context.Context, gameID uuid.UUID) (*HitReport, error)
}

type GameResultAnnouncer interface {
	AnnounceGameResult(chatID int64, report *HitReport)
}

type GameScheduler struct {
	gameRepo        domain.GameRepository
	gunslingerRepo  domain.GunslingerRepository
	userRepo        domain.UserRepository
	chatRepo        domain.ChatRepository
	starter         GameStarter
	announcer       GameResultAnnouncer
	bot             *Bot
	trans           *translator.Translator
	lastMidnightRun string
}

func NewGameScheduler(
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
	userRepo domain.UserRepository,
	chatRepo domain.ChatRepository,
	starter GameStarter,
	announcer GameResultAnnouncer,
	bot *Bot,
	trans *translator.Translator,
) *GameScheduler {
	return &GameScheduler{
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		userRepo:       userRepo,
		chatRepo:       chatRepo,
		starter:        starter,
		announcer:      announcer,
		bot:            bot,
		trans:          trans,
	}
}

func (s *GameScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(SchedulerTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.process(ctx, now)
		}
	}
}

func (s *GameScheduler) process(ctx context.Context, now time.Time) {
	games, err := s.gameRepo.GetActiveDynamicGames()
	if err != nil {
		log.Printf("scheduler: failed to load active dynamic games: %v", err)
		return
	}

	s.processDueGames(ctx, games, now)
	s.processMidnightGames(ctx, games, now)
	s.processFomoGames(ctx, games, now)
	s.processHunger(ctx, now, games)
}

func (s *GameScheduler) processDueGames(ctx context.Context, games []*domain.Game, now time.Time) {
	var dueGames []*domain.Game
	for _, game := range games {
		if !game.StartDeadline.Valid {
			continue
		}

		gunslingers, err := s.gunslingerRepo.GetByGameID(game.ID)
		if err != nil {
			log.Printf("scheduler: failed to load gunslingers for game %s: %v", game.ID, err)
			continue
		}

		if len(gunslingers) >= domain.DynamicMaxPlayers {
			dueGames = append(dueGames, game)
			continue
		}

		if !game.StartDeadline.Time.After(now) {
			dueGames = append(dueGames, game)
		}
	}

	s.startGames(ctx, dueGames)
}

func (s *GameScheduler) processMidnightGames(ctx context.Context, games []*domain.Game, now time.Time) {
	if now.Hour() != SchedulerMidnightHour || now.Minute() != SchedulerMidnightMinute {
		return
	}

	today := now.Format("2006-01-02")
	if s.lastMidnightRun == today {
		return
	}
	s.lastMidnightRun = today

	var midnightGames []*domain.Game
	for _, game := range games {
		if game.StartDeadline.Valid && !game.StartDeadline.Time.After(now) {
			continue
		}

		gunslingers, err := s.gunslingerRepo.GetByGameID(game.ID)
		if err != nil {
			log.Printf("scheduler: failed to load gunslingers for game %s: %v", game.ID, err)
			continue
		}

		if len(gunslingers) >= domain.DynamicMinPlayers {
			midnightGames = append(midnightGames, game)
		}
	}

	s.startGames(ctx, midnightGames)
}

func (s *GameScheduler) startGames(ctx context.Context, games []*domain.Game) {
	for _, game := range games {
		s.startGame(ctx, game)
		time.Sleep(SchedulerStartInterval)
	}
}

func (s *GameScheduler) processFomoGames(ctx context.Context, games []*domain.Game, now time.Time) {
	for _, game := range games {
		if game.FomoNotified {
			continue
		}
		if !game.StartDeadline.Valid {
			continue
		}
		if !game.StartDeadline.Time.After(now) {
			continue
		}

		count := len(game.Gunslingers)
		if count < domain.DynamicMinPlayers {
			continue
		}

		duration := domain.DynamicDeadlineDuration(count)
		if duration == 0 {
			continue
		}

		halftime := game.StartDeadline.Time.Add(-duration / 2)
		if now.Before(halftime) {
			continue
		}

		s.sendHuntMessage(game)
		game.FomoNotified = true
		if err := s.gameRepo.Update(game); err != nil {
			log.Printf("scheduler: failed to update game fomo status: %v", err)
		}
	}
}

func (s *GameScheduler) sendHuntMessage(game *domain.Game) {
	allPlayerIDs, err := s.gunslingerRepo.GetDistinctPlayerIDsInChat(game.ChatID)
	if err != nil {
		log.Printf("scheduler: failed to get distinct player ids for hunt: %v", err)
		return
	}

	currentPlayers := make(map[int64]bool, len(game.Gunslingers))
	for _, gs := range game.Gunslingers {
		currentPlayers[gs.PlayerID] = true
	}

	var candidates []int64
	for _, pid := range allPlayerIDs {
		if !currentPlayers[pid] {
			candidates = append(candidates, pid)
		}
	}

	if len(candidates) == 0 {
		return
	}

	targetID := candidates[rand.Intn(len(candidates))]
	user, err := s.userRepo.Get(targetID)
	if err != nil {
		log.Printf("scheduler: failed to get user for hunt: %v", err)
		return
	}

	message := s.trans.Get("nagant hunt", translator.Config{
		Args: map[string]string{
			"%missing": user.Mention(),
		},
	})

	s.bot.SendMessage(game.ChatID, message)
}

func (s *GameScheduler) processHunger(ctx context.Context, now time.Time, activeGames []*domain.Game) {
	activeChatIDs := make(map[int64]bool, len(activeGames))
	for _, g := range activeGames {
		activeChatIDs[g.ChatID] = true
	}

	chatIDs, err := s.chatRepo.GetAllChatIDs()
	if err != nil {
		log.Printf("scheduler: failed to get all chat IDs: %v", err)
		return
	}

	cooldownStr := config.GetEnv(config.CooldownMinutes)
	cooldown := config.DefaultCooldown
	if minutes, err := strconv.Atoi(cooldownStr); err == nil {
		cooldown = minutes
	}

	for _, chatID := range chatIDs {
		if activeChatIDs[chatID] {
			continue
		}

		if cooldown > 0 {
			lastGame, err := s.gameRepo.GetLastPlayedForChat(chatID)
			if err != nil {
				continue
			}

			if !lastGame.PlayedAt.Valid {
				continue
			}

			if lastGame.PlayedAt.Time.Add(time.Duration(cooldown) * time.Minute).After(now) {
				continue
			}
		}

		chat, err := s.chatRepo.Get(chatID)
		if err != nil {
			log.Printf("scheduler: failed to get chat %d: %v", chatID, err)
			continue
		}

		if !chat.HungerPostedAt.Valid {
			message := s.trans.Get("nagant hunger stage1", translator.Config{})
			s.bot.SendMessage(chatID, message)

			chat.HungerPostedAt = getSQLNullTime(now)
			chat.HungerStage = domain.HungerStageAwakening
			if err := s.chatRepo.Update(&chat); err != nil {
				log.Printf("scheduler: failed to update chat hunger stage: %v", err)
			}
		} else if chat.HungerStage == domain.HungerStageAwakening && now.Sub(chat.HungerPostedAt.Time) > 4*time.Hour {
			message := s.trans.Get("nagant hunger stage2", translator.Config{})
			s.bot.SendMessage(chatID, message)

			chat.HungerPostedAt = getSQLNullTime(now)
			chat.HungerStage = domain.HungerStageImpatience
			if err := s.chatRepo.Update(&chat); err != nil {
				log.Printf("scheduler: failed to update chat hunger stage: %v", err)
			}
		} else if chat.HungerStage == domain.HungerStageImpatience && now.Sub(chat.HungerPostedAt.Time) > 4*time.Hour {
			message := s.trans.Get("nagant hunger stage3", translator.Config{})
			s.bot.SendMessage(chatID, message)

			chat.HungerPostedAt = getSQLNullTime(now)
			chat.HungerStage = domain.HungerStageRage
			if err := s.chatRepo.Update(&chat); err != nil {
				log.Printf("scheduler: failed to update chat hunger stage: %v", err)
			}
		}
	}
}

func (s *GameScheduler) resetHungerForChat(chatID int64) {
	chat, err := s.chatRepo.Get(chatID)
	if err != nil {
		return
	}

	if !chat.IsHungerActive() {
		return
	}

	chat.ResetHunger()
	if err := s.chatRepo.Update(&chat); err != nil {
		log.Printf("scheduler: failed to reset hunger for chat %d: %v", chatID, err)
	}
}

func (s *GameScheduler) startGame(ctx context.Context, game *domain.Game) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("scheduler: recovered from panic while starting game %s: %v", game.ID, r)
		}
	}()

	report, err := s.starter.Execute(ctx, game.ID)
	if err != nil {
		log.Printf("scheduler: failed to start game %s: %v", game.ID, err)
		return
	}

	s.resetHungerForChat(game.ChatID)
	s.announcer.AnnounceGameResult(game.ChatID, report)
}

func getSQLNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}
