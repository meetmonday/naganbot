package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	GameStatusLobby    GameStatus = "lobby"
	GameStatusStarting GameStatus = "starting"
	GameStatusPlayed   GameStatus = "played"
)

// LobbyIdleTTL is how long a game may stay unstarted without anyone joining
// before the scheduler drops it. Without it a single player in a lobby keeps
// HasActiveInChat true forever and the chat can never start another game.
const LobbyIdleTTL = 6 * time.Hour

type Game struct {
	ID            uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	ChatID        int64
	Chat          Chat
	OwnerID       int64
	Owner         User
	Gunslingers   []*Gunslinger
	CreatedAt     time.Time
	PlayedAt      sql.NullTime
	BulletType    string
	ProofURL      string
	PlayersCount  int        `gorm:"->;<-:create;not null;default:6"`
	Mode          GameMode   `gorm:"not null;default:'classic'"`
	Status        GameStatus `gorm:"not null;default:'lobby'"`
	StartDeadline sql.NullTime
}

type GameRepository interface {
	GetByID(uuid.UUID) (*Game, error)
	GetLatestGamesInChat(int64, int) ([]Game, error)
	GetLatestForChat(int64) (*Game, error)
	GetLastPlayedForChat(int64) (*Game, error)
	GetActiveForChat(int64) (*Game, error)
	Store(*Game) error
	Update(*Game) error
	HasActiveInChat(id int64) bool
	GetActiveDynamicGames() ([]*Game, error)
	GetIdleLobbies(before time.Time) ([]*Game, error)
	Delete(*Game) error
}

func NewGame(chatID int64, ownerID int64, playersCount int, mode GameMode) *Game {
	ID := uuid.Must(uuid.NewV7())
	gunslinger := NewGunslinger(ID, ownerID)

	game := &Game{
		ID:           ID,
		ChatID:       chatID,
		OwnerID:      ownerID,
		Gunslingers:  []*Gunslinger{gunslinger},
		CreatedAt:    time.Now(),
		PlayedAt:     sql.NullTime{},
		PlayersCount: playersCount,
		Mode:         mode,
		Status:       GameStatusLobby,
	}
	gunslinger.Game = game

	return game
}

func (g *Game) IsPlayed() bool {
	return g.PlayedAt.Valid
}

func (g *Game) MarkAsPlayed(withBullet string, proofURL string) {
	if !g.IsPlayed() {
		g.PlayedAt = sql.NullTime{Time: time.Now(), Valid: true}
		g.Status = GameStatusPlayed
		g.BulletType = withBullet
		g.ProofURL = proofURL
	}
}

func (g *Game) ShouldStartNow(count int) bool {
	if g.Mode == GameModeClassic {
		return count >= g.PlayersCount
	}

	return count >= DynamicMaxPlayers
}

// LastActivityAt reports when the game was last touched: a join extends the
// life of a lobby, creation is the fallback for a game nobody joined.
func (g *Game) LastActivityAt() time.Time {
	last := g.CreatedAt
	for _, gunslinger := range g.Gunslingers {
		if gunslinger.JoinedAt.After(last) {
			last = gunslinger.JoinedAt
		}
	}

	return last
}

// IsIdle tells whether the game has been waiting for players for longer than
// the cutoff and nobody joined since.
func (g *Game) IsIdle(cutoff time.Time) bool {
	if g.IsPlayed() || g.Status != GameStatusLobby {
		return false
	}

	return g.LastActivityAt().Before(cutoff)
}
