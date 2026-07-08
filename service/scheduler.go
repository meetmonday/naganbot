package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
	starter GameStarter,
	announcer GameResultAnnouncer,
	bot *Bot,
	trans *translator.Translator,
) *GameScheduler {
	return &GameScheduler{
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		userRepo:       userRepo,
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

		totalDuration := game.StartDeadline.Time.Sub(game.CreatedAt)
		halftime := game.CreatedAt.Add(totalDuration / 2)

		if now.Before(halftime) {
			continue
		}

		players, err := s.gunslingerRepo.GetPlayersWithStreakInChat(game.ChatID, game.ID)
		if err != nil {
			log.Printf("scheduler: failed to get players with streaks: %v", err)
			continue
		}

		if len(players) > 0 {
			var playerIDs []int64
			for _, p := range players {
				playerIDs = append(playerIDs, p.PlayerID)
			}

			users, err := s.userRepo.GetByIDs(playerIDs)
			if err != nil {
				log.Printf("scheduler: failed to get users for fomo: %v", err)
				continue
			}

			mentions := make([]string, len(users))
			for i, u := range users {
				mentions[i] = u.Mention()
			}

			message := s.trans.Get("fomo reminder", translator.Config{
				Args: map[string]string{
					"%players": strings.Join(mentions, ", "),
				},
			})

			s.bot.SendMessage(game.ChatID, message)
		}

		game.FomoNotified = true
		if err := s.gameRepo.Update(game); err != nil {
			log.Printf("scheduler: failed to update game fomo status: %v", err)
		}
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

	s.announcer.AnnounceGameResult(game.ChatID, report)
}
