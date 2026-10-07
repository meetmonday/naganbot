package domain

import (
	"testing"
	"time"
)

func TestGame_MarkAsPlayed(t *testing.T) {
	var game *Game

	game = NewGame(0, 0, 6, GameModeClassic)
	game.MarkAsPlayed("lead", "https://example.com/proof")
	if game.IsPlayed() != true {
		t.Error()
	}
}

func TestGame_IsIdle(t *testing.T) {
	now := time.Date(2026, 7, 1, 14, 0, 0, 0, time.UTC)
	cutoff := now.Add(-LobbyIdleTTL)

	tests := []struct {
		name     string
		game     *Game
		expected bool
	}{
		{
			name:     "nobody joined for longer than ttl",
			game:     &Game{Status: GameStatusLobby, CreatedAt: cutoff.Add(-time.Hour)},
			expected: true,
		},
		{
			name:     "recent join extends the life of a lobby",
			game:     &Game{Status: GameStatusLobby, CreatedAt: cutoff.Add(-time.Hour), Gunslingers: []*Gunslinger{{JoinedAt: cutoff.Add(time.Hour)}}},
			expected: false,
		},
		{
			name:     "created recently",
			game:     &Game{Status: GameStatusLobby, CreatedAt: now},
			expected: false,
		},
		{
			name:     "game with a deadline is never idle",
			game:     &Game{Status: GameStatusStarting, CreatedAt: cutoff.Add(-time.Hour)},
			expected: false,
		},
		{
			name:     "played game is never idle",
			game:     &Game{Status: GameStatusPlayed, CreatedAt: cutoff.Add(-time.Hour)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.game.IsIdle(cutoff); got != tt.expected {
				t.Errorf("IsIdle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestGame_LastActivityAt(t *testing.T) {
	now := time.Date(2026, 7, 1, 14, 0, 0, 0, time.UTC)

	game := &Game{
		CreatedAt: now.Add(-3 * time.Hour),
		Gunslingers: []*Gunslinger{
			{JoinedAt: now.Add(-2 * time.Hour)},
			{JoinedAt: now.Add(-time.Minute)},
			{JoinedAt: now.Add(-time.Hour)},
		},
	}

	if got := game.LastActivityAt(); !got.Equal(now.Add(-time.Minute)) {
		t.Errorf("LastActivityAt() = %v, expected %v", got, now.Add(-time.Minute))
	}

	empty := &Game{CreatedAt: now.Add(-time.Hour)}
	if got := empty.LastActivityAt(); !got.Equal(now.Add(-time.Hour)) {
		t.Errorf("LastActivityAt() without gunslingers = %v, expected %v", got, now.Add(-time.Hour))
	}
}

func TestGame_UpdateStartDeadline(t *testing.T) {
	now := time.Date(2026, 7, 1, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		count          int
		expectedStatus GameStatus
		expectedValid  bool
	}{
		{"below minimum", 2, GameStatusLobby, false},
		{"minimum", 3, GameStatusStarting, true},
		{"four players", 4, GameStatusStarting, true},
		{"five players", 5, GameStatusStarting, true},
		{"maximum", 10, GameStatusStarting, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(0, 0, 6, GameModeDynamic)
			game.UpdateStartDeadline(tt.count, now)

			if game.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, game.Status)
			}

			if game.StartDeadline.Valid != tt.expectedValid {
				t.Errorf("expected deadline validity %v, got %v", tt.expectedValid, game.StartDeadline.Valid)
			}
		})
	}
}

func TestGame_UpdateStartDeadline_Midnight(t *testing.T) {
	now := time.Date(2026, 7, 1, 23, 50, 0, 0, time.UTC)
	game := NewGame(0, 0, 6, GameModeDynamic)
	game.UpdateStartDeadline(3, now)

	expected := time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)
	if !game.StartDeadline.Time.Equal(expected) {
		t.Errorf("expected deadline %v, got %v", expected, game.StartDeadline.Time)
	}
}

func TestDynamicDeadlineDuration(t *testing.T) {
	tests := []struct {
		count    int
		expected time.Duration
	}{
		{2, 0},
		{3, 60 * time.Minute},
		{4, 30 * time.Minute},
		{5, 15 * time.Minute},
		{6, 15 * time.Minute},
		{10, 15 * time.Minute},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := DynamicDeadlineDuration(tt.count); got != tt.expected {
				t.Errorf("DynamicDeadlineDuration(%d) = %v, expected %v", tt.count, got, tt.expected)
			}
		})
	}
}

func TestNextMidnight(t *testing.T) {
	tests := []struct {
		name     string
		t        time.Time
		expected time.Time
	}{
		{
			name:     "afternoon",
			t:        time.Date(2026, 7, 1, 14, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "just after midnight",
			t:        time.Date(2026, 7, 1, 0, 0, 1, 0, time.UTC),
			expected: time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextMidnight(tt.t); !got.Equal(tt.expected) {
				t.Errorf("NextMidnight() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
