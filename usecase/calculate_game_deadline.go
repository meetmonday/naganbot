package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
)

type CalculateGameDeadlineUseCase struct {
	gameRepo       domain.GameRepository
	gunslingerRepo domain.GunslingerRepository
}

func NewCalculateGameDeadlineUseCase(
	gameRepo domain.GameRepository,
	gunslingerRepo domain.GunslingerRepository,
) *CalculateGameDeadlineUseCase {
	return &CalculateGameDeadlineUseCase{
		gameRepo:       gameRepo,
		gunslingerRepo: gunslingerRepo,
	}
}

func (uc *CalculateGameDeadlineUseCase) Execute(gameID uuid.UUID, now time.Time) (*domain.Game, int, error) {
	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, 0, err
	}

	gunslingers, err := uc.gunslingerRepo.GetByGameID(game.ID)
	if err != nil {
		return nil, 0, err
	}

	count := len(gunslingers)
	if game.Mode == domain.GameModeDynamic {
		game.UpdateStartDeadline(count, now)
		if err := uc.gameRepo.Update(game); err != nil {
			return nil, 0, err
		}
	}

	return game, count, nil
}
