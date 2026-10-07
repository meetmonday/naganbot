package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
)

type GameRepository struct {
	domain.GameRepository
	orm *gorm.DB
}

func NewGameRepository(
	orm *gorm.DB,
) domain.GameRepository {
	return &GameRepository{
		orm: orm,
	}
}

func (repo GameRepository) GetByID(id uuid.UUID) (*domain.Game, error) {
	var game domain.Game
	err := repo.orm.Where("id = ?", id).First(&game).Error
	return &game, err
}

func (repo GameRepository) GetLatestGamesInChat(chatID int64, limit int) ([]domain.Game, error) {
	var games []domain.Game
	err := repo.getQueryByChat(chatID).
		Where("played_at IS NOT NULL").
		Order("played_at DESC").
		Limit(limit).
		Find(&games).
		Error

	return games, err
}

func (repo GameRepository) GetLatestForChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Order("created_at DESC").
		First(&game).
		Error
	return &game, err
}

func (repo GameRepository) GetLastPlayedForChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Where("played_at IS NOT NULL").
		Order("played_at DESC").
		First(&game).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &game, err
}

func (repo GameRepository) GetActiveForChat(chatID int64) (*domain.Game, error) {
	var game domain.Game
	err := repo.getQueryByChat(chatID).
		Preload("Owner").
		Where("played_at IS NULL").
		First(&game).
		Error
	return &game, err
}

func (repo GameRepository) Store(game *domain.Game) error {
	return repo.orm.Create(game).Error
}

func (repo GameRepository) Update(game *domain.Game) error {
	return repo.orm.Select("Status", "PlayedAt", "BulletType", "ProofURL", "StartDeadline", "PlayersCount").Updates(game).Error
}

func (repo GameRepository) HasActiveInChat(chatID int64) bool {
	var counter int64
	repo.orm.Model(&domain.Game{}).
		Where("chat_id = ?", chatID).
		Where("played_at IS NULL").
		Count(&counter)

	return counter > 0
}

// GetIdleLobbies returns unplayed games that nobody joined recently. The coarse
// filter runs in SQL, the precise one in domain.Game.IsIdle.
func (repo GameRepository) GetIdleLobbies(cutoff time.Time) ([]*domain.Game, error) {
	var games []*domain.Game
	err := repo.orm.
		Preload("Gunslingers", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_at ASC")
		}).
		Where("played_at IS NULL").
		Where("status = ?", domain.GameStatusLobby).
		Where("created_at < ?", cutoff).
		Find(&games).
		Error
	if err != nil {
		return nil, err
	}

	idle := make([]*domain.Game, 0, len(games))
	for _, game := range games {
		if game.IsIdle(cutoff) {
			idle = append(idle, game)
		}
	}

	return idle, nil
}

// Delete removes the game together with its gunslingers, otherwise the foreign
// key keeps the row alive.
func (repo GameRepository) Delete(game *domain.Game) error {
	return repo.orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("game_id = ?", game.ID).Delete(&domain.Gunslinger{}).Error; err != nil {
			return err
		}

		return tx.Where("id = ?", game.ID).Delete(&domain.Game{}).Error
	})
}

func (repo GameRepository) GetActiveDynamicGames() ([]*domain.Game, error) {
	var games []*domain.Game
	err := repo.orm.
		Preload("Gunslingers", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_at ASC")
		}).
		Preload("Gunslingers.Player").
		Preload("Owner").
		Preload("Chat").
		Where("played_at IS NULL").
		Where("mode = ?", domain.GameModeDynamic).
		Find(&games).
		Error

	return games, err
}

func (repo GameRepository) getQueryByChat(chatID int64) *gorm.DB {
	return repo.orm.
		Preload("Gunslingers", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_at ASC")
		}).
		Preload("Gunslingers.Player").
		Preload("Owner").
		Where("chat_id = ?", chatID)
}
