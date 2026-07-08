package service

import (
	"time"

	"github.com/taranovegor/naganbot/translator"
)

type GameAnnouncer struct {
	bot   *Bot
	trans *translator.Translator
}

func NewGameAnnouncer(
	bot *Bot,
	trans *translator.Translator,
) *GameAnnouncer {
	return &GameAnnouncer{
		bot:   bot,
		trans: trans,
	}
}

func (a *GameAnnouncer) AnnounceGameResult(chatID int64, report *HitReport) {
	for _, message := range a.trans.GetMany("play the game", translator.Config{}) {
		a.bot.SendMessage(chatID, message)
		time.Sleep(time.Second)
	}

	isAtomic := report.BulletType == BulletAtomicType

	if isAtomic {
		a.bot.SendMessage(chatID, a.trans.Get("killed by atomic bullet", translator.Config{}))
	}

	for _, victim := range report.Victims {
		if !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
		}

		err := a.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("player is not kicked", translator.Config{}))
		}
	}
}
