package domain

import (
	"database/sql"
	"time"
)

const (
	DynamicMinPlayers = 3
	DynamicMaxPlayers = 10
)

var DynamicDeadlines = map[int]time.Duration{
	3: 60 * time.Minute,
	4: 30 * time.Minute,
	5: 15 * time.Minute,
}

func DynamicDeadlineDuration(count int) time.Duration {
	if count < DynamicMinPlayers {
		return 0
	}

	for threshold := 5; threshold >= 3; threshold-- {
		if count >= threshold {
			return DynamicDeadlines[threshold]
		}
	}

	return 0
}

func NextMidnight(t time.Time) time.Time {
	y, m, d := t.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	if !t.Before(midnight) {
		midnight = midnight.Add(24 * time.Hour)
	}

	return midnight
}

func (g *Game) UpdateStartDeadline(count int, now time.Time) {
	if g.Mode != GameModeDynamic || count < DynamicMinPlayers {
		g.Status = GameStatusLobby
		g.StartDeadline = sql.NullTime{}
		return
	}

	if count >= DynamicMaxPlayers {
		g.Status = GameStatusStarting
		g.StartDeadline = sql.NullTime{Time: now, Valid: true}
		return
	}

	dynamicDeadline := now.Add(DynamicDeadlineDuration(count))
	midnight := NextMidnight(now)

	deadline := dynamicDeadline
	if midnight.Before(deadline) {
		deadline = midnight
	}

	g.Status = GameStatusStarting
	g.StartDeadline = sql.NullTime{Time: deadline, Valid: true}
}
