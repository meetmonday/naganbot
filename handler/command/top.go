package command

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type TopHandler struct {
	Handler
	bot        *service.Bot
	trans      *translator.Translator
	user       domain.UserRepository
	gunslinger domain.GunslingerRepository
}

func NewTopHandler(
	bot *service.Bot,
	trans *translator.Translator,
	user domain.UserRepository,
	gunslinger domain.GunslingerRepository,
) Handler {
	return &TopHandler{
		bot:        bot,
		trans:      trans,
		user:       user,
		gunslinger: gunslinger,
	}
}

func (hdlr TopHandler) Name() string {
	return "top"
}

func (hdlr TopHandler) Execute(msg *tgbotapi.Message) {
	var players []domain.GunslingerTopShotPlayer
	var err error
	chatID := msg.Chat.ID
	if year, parseErr := strconv.Atoi(msg.CommandArguments()); parseErr == nil {
		players, err = hdlr.gunslinger.GetTopShotPlayersByYearInChat(chatID, year)
	} else {
		players, err = hdlr.gunslinger.GetTopShotPlayersInChat(chatID)
	}
	if err != nil {
		log.Printf("failed to get top players: %v", err)
		hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
		return
	}

	var message string
	if len(players) == 0 {
		message = hdlr.trans.Get("top is not determined", translator.Config{})
	} else {
		ids := make([]int64, len(players))
		for i, p := range players {
			ids[i] = p.PlayerId
		}
		users, err := hdlr.user.GetByIDs(ids)
		if err != nil {
			log.Printf("failed to get users: %v", err)
			hdlr.bot.SendMessage(chatID, hdlr.trans.Get("something went wrong", translator.Config{}))
			return
		}
		userMap := make(map[int64]domain.User, len(users))
		for _, u := range users {
			userMap[u.ID] = u
		}

		message = hdlr.trans.Get("top players by games", translator.Config{
			Args: map[string]string{"%number": strconv.Itoa(len(players))},
		})
		for i, player := range players {
			user, ok := userMap[player.PlayerId]
			if !ok {
				continue
			}

			message += "\n" + hdlr.trans.Get("top game player", translator.Config{
				Args: map[string]string{
					"%i":     strconv.Itoa(i + 1),
					"%user":  user.Name(),
					"%times": strconv.Itoa(player.Times),
				},
				Count: player.Times,
			})
		}
	}

	keyboard := hdlr.buildKeyboard("shot")
	hdlr.bot.SendInlineKeyboard(chatID, message, keyboard)
}

func (hdlr TopHandler) buildKeyboard(activeTab string) service.InlineKeyboard {
	tabs := []struct {
		key      string
		transKey string
	}{
		{"shot", "top_tab_shot"},
		{"creators", "top_tab_creators"},
		{"active", "top_tab_active"},
		{"streak", "top_tab_streak"},
	}

	var row1, row2 []service.KeyboardButton
	for i, t := range tabs {
		label := hdlr.trans.Get(t.transKey, translator.Config{})
		if t.key == activeTab {
			label = "▸ " + label
		}
		btn := service.KeyboardButton{
			Data:  callback.TopTab.SetArgs(t.key).ToString(),
			Label: label,
		}
		if i < 2 {
			row1 = append(row1, btn)
		} else {
			row2 = append(row2, btn)
		}
	}

	return service.InlineKeyboard{row1, row2}
}
