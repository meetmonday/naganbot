package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"gorm.io/gorm"
	"sort"
)

type GunslingerRepository struct {
	domain.GunslingerRepository
	orm *gorm.DB
}

func NewGunslingerRepository(
	orm *gorm.DB,
) domain.GunslingerRepository {
	return &GunslingerRepository{
		orm: orm,
	}
}

func (repo GunslingerRepository) Store(gunslinger *domain.Gunslinger) error {
	return repo.orm.Create(gunslinger).Error
}

func (repo GunslingerRepository) Update(gunslingers []*domain.Gunslinger) error {
	return repo.orm.Transaction(func(tx *gorm.DB) error {
		for _, gunslinger := range gunslingers {
			if err := tx.Updates(gunslinger).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo GunslingerRepository) GetByGameID(gameID uuid.UUID) ([]*domain.Gunslinger, error) {
	var gunslingers []*domain.Gunslinger
	err := repo.orm.
		Preload("Game").
		Preload("Player").
		Where("game_id = ?", gameID).
		Order("joined_at ASC").
		Find(&gunslingers).
		Error

	return gunslingers, err
}

func (repo GunslingerRepository) IsPlayerExistsInGame(userID int64, gameID uuid.UUID) bool {
	var counter int64
	repo.orm.Model(&domain.Gunslinger{}).
		Where("player_id = ?", userID).
		Where("game_id = ?", gameID).
		Count(&counter)

	return counter > 0
}

func (repo GunslingerRepository) GetTopShotPlayersInChat(chatID int64) ([]domain.GunslingerTopShotPlayer, error) {
	var players []domain.GunslingerTopShotPlayer
	err := repo.getQueryTopShotPlayersInChat(chatID).
		Find(&players).
		Error

	return players, err
}

func (repo GunslingerRepository) GetTopShotPlayersByYearInChat(chatID int64, year int) ([]domain.GunslingerTopShotPlayer, error) {
	var players []domain.GunslingerTopShotPlayer
	err := repo.getQueryTopShotPlayersInChat(chatID).
		Where("EXTRACT(YEAR FROM games.created_at) = ?", year).
		Find(&players).
		Error

	return players, err
}

func (repo GunslingerRepository) CountNumberOfPlayerGamesInChat(userID int64, chatID int64) int64 {
	var counter int64
	repo.getQueryPlayerGamesInChat(userID, chatID).
		Count(&counter)

	return counter
}

func (repo GunslingerRepository) CountNumberOfSelfShotsInChat(userID int64, chatID int64) int64 {
	var counter int64
	repo.getQueryPlayerGamesInChat(userID, chatID).
		Where("shot_himself = ?", true).
		Count(&counter)

	return counter
}

func (repo GunslingerRepository) GetPlayerStreaks(userID int64, chatID int64) (int, int) {
	type row struct {
		Participated int
		ShotHimself  int
	}

	rows, err := repo.orm.Raw(`
		SELECT
			CASE WHEN gs.player_id IS NOT NULL THEN 1 ELSE 0 END AS participated,
			CASE WHEN gs.shot_himself IS TRUE THEN 1 ELSE 0 END AS shot_himself
		FROM games g
		LEFT JOIN gunslingers gs ON gs.game_id = g.id AND gs.player_id = ?
		WHERE g.chat_id = ? AND g.played_at IS NOT NULL
		ORDER BY g.played_at DESC
	`, userID, chatID).Rows()
	if err != nil {
		return 0, 0
	}
	defer rows.Close()

	participationStreak := 0
	lossStreak := 0

	for rows.Next() {
		var r row
		if err := rows.Scan(&r.Participated, &r.ShotHimself); err != nil {
			break
		}

		if r.Participated == 0 {
			break
		}

		participationStreak++
		if r.ShotHimself == 1 {
			lossStreak++
		} else {
			lossStreak = 0
		}
	}

	return participationStreak, lossStreak
}

func (repo GunslingerRepository) GetPlayersWithStreakInChat(chatID int64, excludeGameID uuid.UUID) ([]domain.PlayerStreak, error) {
	type gamePlayerPair struct {
		GameID   uuid.UUID
		PlayerID int64
	}

	var pairs []gamePlayerPair
	err := repo.orm.Table("games").
		Select("games.id AS game_id, gunslingers.player_id").
		Joins("INNER JOIN gunslingers ON gunslingers.game_id = games.id").
		Where("games.chat_id = ?", chatID).
		Where("games.played_at IS NOT NULL").
		Order("games.played_at DESC").
		Limit(100).
		Find(&pairs).Error
	if err != nil {
		return nil, err
	}

	excludePlayers := make(map[int64]bool)
	var exPlayerIDs []int64
	repo.orm.Table("gunslingers").
		Select("player_id").
		Where("game_id = ?", excludeGameID).
		Find(&exPlayerIDs)
	for _, pid := range exPlayerIDs {
		excludePlayers[pid] = true
	}

	gamePlayersMap := make(map[uuid.UUID][]int64)
	var gameOrder []uuid.UUID
	for _, p := range pairs {
		players, ok := gamePlayersMap[p.GameID]
		if !ok {
			gameOrder = append(gameOrder, p.GameID)
		}
		gamePlayersMap[p.GameID] = append(players, p.PlayerID)
	}

	playerStreak := make(map[int64]int)
	for _, gid := range gameOrder {
		players := gamePlayersMap[gid]
		playerSet := make(map[int64]bool, len(players))
		for _, pid := range players {
			playerSet[pid] = true
			playerStreak[pid]++
		}

		for pid := range playerStreak {
			if !playerSet[pid] {
				delete(playerStreak, pid)
			}
		}
	}

	var result []domain.PlayerStreak
	for pid, streak := range playerStreak {
		if streak > 0 && !excludePlayers[pid] {
			result = append(result, domain.PlayerStreak{
				PlayerID:            pid,
				ParticipationStreak: streak,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ParticipationStreak > result[j].ParticipationStreak
	})

	return result, nil
}

func (repo GunslingerRepository) getQueryTopShotPlayersInChat(chatID int64) *gorm.DB {
	return repo.orm.Table("gunslingers").
		Select("player_id, COUNT(chat_id) as times").
		Joins("INNER JOIN games ON gunslingers.game_id = games.id").
		Where("games.chat_id = ?", chatID).
		Where("games.played_at IS NOT NULL").
		Where("shot_himself = ?", true).
		Group("player_id").
		Order("times DESC").
		Limit(10)
}

func (repo GunslingerRepository) getQueryPlayerGamesInChat(userID int64, chatID int64) *gorm.DB {
	return repo.orm.Model(&domain.Gunslinger{}).
		InnerJoins("Game").
		Where("player_id = ?", userID).
		Where("chat_id = ?", chatID).
		Where("played_at IS NOT NULL")
}
