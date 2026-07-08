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

	a.bot.SendMessage(chatID, a.trans.Get("nagant chosen", translator.Config{}))
	time.Sleep(time.Second)

	isAtomic := report.BulletType == BulletAtomicType

	if isAtomic {
		a.bot.SendMessage(chatID, a.trans.Get("killed by atomic bullet", translator.Config{}))
		time.Sleep(time.Millisecond * 500)
	}

	for _, victim := range report.Victims {
		if !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("gunslinger killed", translator.Config{
				Args: map[string]string{"%gunslinger": victim.Player.Mention()},
			}))
			time.Sleep(time.Millisecond * 500)
		}

		a.bot.SendMessage(chatID, a.trans.Get("nagant final words victim", translator.Config{
			Args: map[string]string{"%user": victim.Player.Mention()},
		}))

		err := a.bot.Kick(chatID, victim.PlayerID)
		if err != nil && !isAtomic {
			a.bot.SendMessage(chatID, a.trans.Get("player is not kicked", translator.Config{}))
		}
	}

	if !isAtomic && len(report.Victims) == 1 {
		a.bot.SendMessage(chatID, a.trans.Get("nagant final words survivors", translator.Config{}))
	}
}

func (a *GameAnnouncer) AnnounceDelayedGameStart(chatID int64, hungerStage int) {
	if hungerStage < 1 || hungerStage > 3 {
		return
	}

	a.bot.SendMessage(chatID, a.trans.Get("nagant game start delayed", translator.Config{}))
	time.Sleep(time.Second)
}
