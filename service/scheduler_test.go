package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
)

type mockGameRepoForScheduler struct {
	games        []*domain.Game
	idleLobbies  []*domain.Game
	deletedGames []uuid.UUID
}

func (m *mockGameRepoForScheduler) GetByID(id uuid.UUID) (*domain.Game, error) {
	for _, g := range m.games {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, nil
}

func (m *mockGameRepoForScheduler) GetLatestGamesInChat(chatID int64, limit int) ([]domain.Game, error) {
	return nil, nil
}

func (m *mockGameRepoForScheduler) GetLatestForChat(chatID int64) (*domain.Game, error) {
	return nil, nil
}

func (m *mockGameRepoForScheduler) GetLastPlayedForChat(chatID int64) (*domain.Game, error) {
	return nil, nil
}

func (m *mockGameRepoForScheduler) GetActiveForChat(chatID int64) (*domain.Game, error) {
	return nil, nil
}

func (m *mockGameRepoForScheduler) Store(game *domain.Game) error {
	return nil
}

func (m *mockGameRepoForScheduler) Update(game *domain.Game) error {
	return nil
}

func (m *mockGameRepoForScheduler) HasActiveInChat(id int64) bool {
	return false
}

func (m *mockGameRepoForScheduler) GetActiveDynamicGames() ([]*domain.Game, error) {
	return m.games, nil
}

// GetIdleLobbies mirrors the repository: it applies domain.Game.IsIdle to the
// candidates so the scheduler is only ever handed genuinely idle games.
func (m *mockGameRepoForScheduler) GetIdleLobbies(cutoff time.Time) ([]*domain.Game, error) {
	idle := make([]*domain.Game, 0, len(m.idleLobbies))
	for _, game := range m.idleLobbies {
		if game.IsIdle(cutoff) {
			idle = append(idle, game)
		}
	}

	return idle, nil
}

func (m *mockGameRepoForScheduler) Delete(game *domain.Game) error {
	m.deletedGames = append(m.deletedGames, game.ID)
	return nil
}

type mockGunslingerRepoForScheduler struct {
	byGame map[uuid.UUID][]*domain.Gunslinger
}

func (m *mockGunslingerRepoForScheduler) Store(gunslinger *domain.Gunslinger) error {
	return nil
}

func (m *mockGunslingerRepoForScheduler) Update(gunslingers []*domain.Gunslinger) error {
	return nil
}

func (m *mockGunslingerRepoForScheduler) GetByGameID(gameID uuid.UUID) ([]*domain.Gunslinger, error) {
	return m.byGame[gameID], nil
}

func (m *mockGunslingerRepoForScheduler) IsPlayerExistsInGame(userID int64, gameID uuid.UUID) bool {
	return false
}

func (m *mockGunslingerRepoForScheduler) GetTopShotPlayersInChat(chatID int64) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}

func (m *mockGunslingerRepoForScheduler) GetTopShotPlayersByYearInChat(chatID int64, year int) ([]domain.GunslingerTopShotPlayer, error) {
	return nil, nil
}

func (m *mockGunslingerRepoForScheduler) CountNumberOfPlayerGamesInChat(userID int64, chatID int64) int64 {
	return 0
}

func (m *mockGunslingerRepoForScheduler) CountNumberOfSelfShotsInChat(userID int64, chatID int64) int64 {
	return 0
}

type mockStarter struct {
	started []uuid.UUID
}

func (m *mockStarter) Execute(ctx context.Context, gameID uuid.UUID) (*HitReport, error) {
	m.started = append(m.started, gameID)
	return &HitReport{BulletType: BulletLeadType}, nil
}

type mockAnnouncer struct {
	announced []int64
}

func (m *mockAnnouncer) AnnounceGameResult(chatID int64, report *HitReport) {
	m.announced = append(m.announced, chatID)
}

func newScheduler(
	games []*domain.Game,
	gunslingers map[uuid.UUID][]*domain.Gunslinger,
) (*GameScheduler, *mockStarter, *mockAnnouncer) {
	return newSchedulerWithLobbies(games, gunslingers, nil)
}

func newSchedulerWithLobbies(
	games []*domain.Game,
	gunslingers map[uuid.UUID][]*domain.Gunslinger,
	idleLobbies []*domain.Game,
) (*GameScheduler, *mockStarter, *mockAnnouncer) {
	starter := &mockStarter{}
	announcer := &mockAnnouncer{}
	scheduler := NewGameScheduler(
		&mockGameRepoForScheduler{games: games, idleLobbies: idleLobbies},
		&mockGunslingerRepoForScheduler{byGame: gunslingers},
		starter,
		announcer,
	)

	return scheduler, starter, announcer
}

func TestGameScheduler_ProcessDueGames(t *testing.T) {
	dueGame := &domain.Game{
		ID:            uuid.Must(uuid.NewV7()),
		ChatID:        1,
		Mode:          domain.GameModeDynamic,
		StartDeadline: sql.NullTime{Time: time.Now().Add(-time.Minute), Valid: true},
	}
	futureGame := &domain.Game{
		ID:            uuid.Must(uuid.NewV7()),
		ChatID:        2,
		Mode:          domain.GameModeDynamic,
		StartDeadline: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
	}

	gunslingers := map[uuid.UUID][]*domain.Gunslinger{
		dueGame.ID:    {{}, {}, {}},
		futureGame.ID: {{}, {}, {}},
	}

	scheduler, starter, announcer := newScheduler([]*domain.Game{dueGame, futureGame}, gunslingers)
	scheduler.process(context.Background(), time.Now())

	if len(starter.started) != 1 || starter.started[0] != dueGame.ID {
		t.Errorf("expected due game to be started, got %v", starter.started)
	}

	if len(announcer.announced) != 1 || announcer.announced[0] != dueGame.ChatID {
		t.Errorf("expected due game to be announced, got %v", announcer.announced)
	}
}

func TestGameScheduler_ProcessMidnightGames(t *testing.T) {
	midnight := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	futureGame := &domain.Game{
		ID:            uuid.Must(uuid.NewV7()),
		ChatID:        1,
		Mode:          domain.GameModeDynamic,
		StartDeadline: sql.NullTime{Time: midnight.Add(time.Hour), Valid: true},
	}

	gunslingers := map[uuid.UUID][]*domain.Gunslinger{
		futureGame.ID: {{}, {}, {}},
	}

	scheduler, starter, announcer := newScheduler([]*domain.Game{futureGame}, gunslingers)
	scheduler.process(context.Background(), midnight)

	if len(starter.started) != 1 || starter.started[0] != futureGame.ID {
		t.Errorf("expected midnight game to be started, got %v", starter.started)
	}

	if len(announcer.announced) != 1 || announcer.announced[0] != futureGame.ChatID {
		t.Errorf("expected midnight game to be announced, got %v", announcer.announced)
	}
}

func TestGameScheduler_ProcessMidnightGames_OnlyOncePerDay(t *testing.T) {
	midnight := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	game := &domain.Game{
		ID:            uuid.Must(uuid.NewV7()),
		ChatID:        1,
		Mode:          domain.GameModeDynamic,
		StartDeadline: sql.NullTime{Time: midnight.Add(time.Hour), Valid: true},
	}

	gunslingers := map[uuid.UUID][]*domain.Gunslinger{
		game.ID: {{}, {}, {}},
	}

	scheduler, starter, _ := newScheduler([]*domain.Game{game}, gunslingers)
	scheduler.process(context.Background(), midnight)
	scheduler.process(context.Background(), midnight.Add(time.Minute))

	if len(starter.started) != 1 {
		t.Errorf("expected midnight game to be started only once, got %v", starter.started)
	}
}

func TestGameScheduler_DiscardIdleLobbies(t *testing.T) {
	idleGame := &domain.Game{
		ID:        uuid.Must(uuid.NewV7()),
		ChatID:    1,
		Mode:      domain.GameModeDynamic,
		Status:    domain.GameStatusLobby,
		CreatedAt: time.Now().Add(-domain.LobbyIdleTTL - time.Hour),
	}

	scheduler, _, _ := newSchedulerWithLobbies(nil, nil, []*domain.Game{idleGame})
	scheduler.process(context.Background(), time.Now())

	repo := scheduler.gameRepo.(*mockGameRepoForScheduler)
	if len(repo.deletedGames) != 1 || repo.deletedGames[0] != idleGame.ID {
		t.Errorf("expected idle lobby to be discarded, got %v", repo.deletedGames)
	}
}

func TestGameScheduler_DiscardIdleLobbies_KeepsRecentlyJoined(t *testing.T) {
	now := time.Now()
	recent := &domain.Game{
		ID:        uuid.Must(uuid.NewV7()),
		ChatID:    1,
		Mode:      domain.GameModeDynamic,
		Status:    domain.GameStatusLobby,
		CreatedAt: now.Add(-domain.LobbyIdleTTL - time.Hour),
		Gunslingers: []*domain.Gunslinger{
			{JoinedAt: now.Add(-time.Minute)},
		},
	}

	scheduler, _, _ := newSchedulerWithLobbies(nil, nil, []*domain.Game{recent})
	scheduler.process(context.Background(), now)

	repo := scheduler.gameRepo.(*mockGameRepoForScheduler)
	if len(repo.deletedGames) != 0 {
		t.Errorf("expected lobby with a recent join to be kept, got %v", repo.deletedGames)
	}
}

func TestGameScheduler_DiscardIdleLobbies_KeepsStartingGame(t *testing.T) {
	starting := &domain.Game{
		ID:        uuid.Must(uuid.NewV7()),
		ChatID:    1,
		Mode:      domain.GameModeDynamic,
		Status:    domain.GameStatusStarting,
		CreatedAt: time.Now().Add(-domain.LobbyIdleTTL - time.Hour),
	}

	scheduler, _, _ := newSchedulerWithLobbies(nil, nil, []*domain.Game{starting})
	scheduler.process(context.Background(), time.Now())

	repo := scheduler.gameRepo.(*mockGameRepoForScheduler)
	if len(repo.deletedGames) != 0 {
		t.Errorf("expected game with a deadline to be kept, got %v", repo.deletedGames)
	}
}
