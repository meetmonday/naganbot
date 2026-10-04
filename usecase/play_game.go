package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
)

var (
	ErrGameAlreadyPlayed = errors.New("game already played")
	ErrNotEnoughPlayers  = errors.New("not enough players")
)

type PlayGameUseCase struct {
	locker         service.Locker
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
	nagan          *service.Nagan
}

func NewPlayGameUseCase(
	locker service.Locker,
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
	nagan *service.Nagan,
) *PlayGameUseCase {
	return &PlayGameUseCase{
		locker:         locker,
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
		nagan:          nagan,
	}
}

func (uc *PlayGameUseCase) Execute(ctx context.Context, gameID uuid.UUID) (*service.HitReport, error) {
	locker := uc.locker.LockFor(fmt.Sprintf("play-game-%d", gameID.ID()))
	if !locker.TryLock() {
		return nil, service.ErrLockFailed
	}
	defer locker.Unlock()

	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	if game.IsPlayed() {
		return nil, ErrGameAlreadyPlayed
	}

	gunslingers, err := uc.gunslingerRepo.GetByGameID(game.ID)
	if err != nil {
		return nil, err
	}

	if !uc.canStart(game, gunslingers, time.Now()) {
		return nil, ErrNotEnoughPlayers
	}

	return uc.play(ctx, game, gunslingers)
}

func (uc *PlayGameUseCase) canStart(game *domain.Game, gunslingers []*domain.Gunslinger, now time.Time) bool {
	if game.Mode == domain.GameModeClassic {
		return len(gunslingers) >= game.PlayersCount
	}

	count := len(gunslingers)
	if count < domain.DynamicMinPlayers {
		return false
	}

	if count >= domain.DynamicMaxPlayers {
		return true
	}

	return game.StartDeadline.Valid && !game.StartDeadline.Time.After(now)
}

func (uc *PlayGameUseCase) ForceExecute(ctx context.Context, gameID uuid.UUID) (*service.HitReport, error) {
	locker := uc.locker.LockFor(fmt.Sprintf("play-game-%d", gameID.ID()))
	if !locker.TryLock() {
		return nil, service.ErrLockFailed
	}
	defer locker.Unlock()

	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	if game.IsPlayed() {
		return nil, ErrGameAlreadyPlayed
	}

	gunslingers, err := uc.gunslingerRepo.GetByGameID(game.ID)
	if err != nil {
		return nil, err
	}

	if len(gunslingers) == 0 {
		return nil, ErrNotEnoughPlayers
	}

	return uc.play(ctx, game, gunslingers)
}

func (uc *PlayGameUseCase) play(ctx context.Context, game *domain.Game, gunslingers []*domain.Gunslinger) (*service.HitReport, error) {
	report, err := uc.nagan.Shoot(ctx, game.ID, gunslingers)
	if err != nil {
		return nil, err
	}
	for _, victim := range report.Victims {
		victim.MarkAsShotHimself()
	}

	game.MarkAsPlayed(report.BulletType, report.ProofURL)
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, err
	}

	if err := uc.gunslingerRepo.Update(report.Victims); err != nil {
		return nil, err
	}

	return report, nil
}
