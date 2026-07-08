package domain

import "database/sql"

type GameMode string

const (
	GameModeClassic GameMode = "classic"
	GameModeDynamic GameMode = "dynamic"
)

type Chat struct {
	ID       int64 `gorm:"primary_key;auto_increment:false"`
	Title    sql.NullString
	Username sql.NullString
	Settings struct {
		Mode            GameMode `gorm:"not null;default:'dynamic'"`
		RequiredPlayers int      `gorm:"not null;default:6"`
	} `gorm:"embedded"`
}

type ChatRepository interface {
	Exists(int64) bool
	Get(int64) (Chat, error)
	Store(*Chat) error
	Update(*Chat) error
}

func NewChat(id int64, title string, username string) *Chat {
	chat := &Chat{ID: id}

	if len(title) > 0 {
		chat.Title = sql.NullString{String: title, Valid: true}
	}

	if len(username) > 0 {
		chat.Username = sql.NullString{String: username, Valid: true}
	}

	return chat
}
