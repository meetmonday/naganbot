package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"strconv"
)

type StatHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	gunslinger domain.GunslingerRepository
}

func NewStatHandler(
	bot *service.Bot,
	trans *translator.Translator,
	gunslinger domain.GunslingerRepository,
) Handler {
	return &StatHandler{
		bot:        bot,
		trans:      trans,
		gunslinger: gunslinger,
	}
}

func (hdlr StatHandler) Name() string {
	return "stat"
}

func (hdlr StatHandler) Execute(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID
	numberOfGames := hdlr.gunslinger.CountNumberOfPlayerGamesInChat(userID, chatID)
	numberOfShotHimself := hdlr.gunslinger.CountNumberOfSelfShotsInChat(userID, chatID)
	participationStreak, lossStreak := hdlr.gunslinger.GetPlayerStreaks(userID, chatID)

	message := hdlr.trans.Get("user game statistics games", translator.Config{
		Args:  map[string]string{"%games": strconv.FormatInt(numberOfGames, 10)},
		Count: int(numberOfGames),
	})
	message += hdlr.trans.Get("user game statistics shots", translator.Config{
		Args:  map[string]string{"%shots": strconv.FormatInt(numberOfShotHimself, 10)},
		Count: int(numberOfShotHimself),
	})
	message += hdlr.trans.Get("user game statistics participation streak", translator.Config{
		Args:  map[string]string{"%ps_games": strconv.Itoa(participationStreak)},
		Count: participationStreak,
	})
	message += hdlr.trans.Get("user game statistics loss streak", translator.Config{
		Args:  map[string]string{"%ls_games": strconv.Itoa(lossStreak)},
		Count: lossStreak,
	})

	hdlr.bot.SendMessage(chatID, message)
}
