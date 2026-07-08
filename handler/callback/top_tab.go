package callback

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
)

type topTab struct {
	bot        *service.Bot
	trans      *translator.Translator
	user       domain.UserRepository
	gunslinger domain.GunslingerRepository
}

func NewTopTab(
	bot *service.Bot,
	trans *translator.Translator,
	user domain.UserRepository,
	gunslinger domain.GunslingerRepository,
) Handler {
	return &topTab{
		bot:        bot,
		trans:      trans,
		user:       user,
		gunslinger: gunslinger,
	}
}

func (h *topTab) Pattern() Pattern {
	return TopTab
}

func (h *topTab) Execute(query *tgbotapi.CallbackQuery) {
	tab := TopTab.GetArg(query.Data, 1)

	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	players, headerKey, itemKey, err := h.getTopData(chatID, tab)
	if err != nil {
		log.Printf("failed to get top data: %v", err)
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}

	if len(players) == 0 {
		err = h.bot.EditMessage(chatID, messageID, h.trans.Get("top is not determined", translator.Config{}), h.buildKeyboard(tab))
		h.bot.AnswerCallback(query.ID, "")
		if err != nil {
			log.Printf("failed to edit message: %v", err)
		}
		return
	}

	ids := make([]int64, len(players))
	for i, p := range players {
		ids[i] = p.PlayerID()
	}
	users, err := h.user.GetByIDs(ids)
	if err != nil {
		log.Printf("failed to get users: %v", err)
		h.bot.AnswerCallback(query.ID, h.trans.Get("something went wrong", translator.Config{}))
		return
	}
	userMap := make(map[int64]domain.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	message := h.trans.Get(headerKey, translator.Config{
		Args: map[string]string{"%number": strconv.Itoa(len(players))},
	})
	for i, player := range players {
		user, ok := userMap[player.PlayerID()]
		if !ok {
			continue
		}

		if ps, ok := player.(streakPlayer); ok {
			key := "top streak player"
			args := map[string]string{
				"%i":     strconv.Itoa(i + 1),
				"%user":  user.Name(),
				"%times": strconv.Itoa(ps.ParticipationStreak),
				"%peak":  strconv.Itoa(ps.PeakStreak),
			}
			if ps.ParticipationStreak == ps.PeakStreak {
				key = "top streak player peak"
				args = map[string]string{
					"%i":     strconv.Itoa(i + 1),
					"%user":  user.Name(),
					"%times": strconv.Itoa(ps.ParticipationStreak),
				}
			}
			streakTxt := h.trans.Get(key, translator.Config{Args: args, Count: ps.ParticipationStreak})
			message += "\n" + streakTxt
		} else {
			message += "\n" + h.trans.Get(itemKey, translator.Config{
				Args: map[string]string{
					"%i":     strconv.Itoa(i + 1),
					"%user":  user.Name(),
					"%times": strconv.Itoa(player.Count()),
				},
				Count: player.Count(),
			})
		}
	}

	err = h.bot.EditMessage(chatID, messageID, message, h.buildKeyboard(tab))
	h.bot.AnswerCallback(query.ID, "")
	if err != nil {
		log.Printf("failed to edit message: %v", err)
	}
}

func (h *topTab) getTopData(chatID int64, tab string) ([]topPlayerItem, string, string, error) {
	const limit = 10

		switch tab {
	case "creators":
		players, err := h.gunslinger.GetTopCreatorsInChat(chatID, limit)
		if err != nil {
			return nil, "", "", err
		}
		items := make([]topPlayerItem, len(players))
		for i, p := range players {
			items[i] = simplePlayer{PID: p.PlayerId, times: p.Times}
		}
		return items, "top creators header", "top game player", nil

	case "active":
		players, err := h.gunslinger.GetTopActivePlayersInChat(chatID, limit)
		if err != nil {
			return nil, "", "", err
		}
		items := make([]topPlayerItem, len(players))
		for i, p := range players {
			items[i] = simplePlayer{PID: p.PlayerId, times: p.Times}
		}
		return items, "top active header", "top game player", nil

	case "streak":
		players, err := h.gunslinger.GetTopStreaksInChat(chatID, limit)
		if err != nil {
			return nil, "", "", err
		}
		items := make([]topPlayerItem, len(players))
		for i, p := range players {
			items[i] = streakPlayer{
				PID:                 p.PlayerID,
				ParticipationStreak: p.ParticipationStreak,
				PeakStreak:          p.PeakStreak,
			}
		}
		return items, "top streak header", "top streak player", nil

	default:
		players, err := h.gunslinger.GetTopShotPlayersInChat(chatID)
		if err != nil {
			return nil, "", "", err
		}
		items := make([]topPlayerItem, len(players))
		for i, p := range players {
			items[i] = simplePlayer{PID: p.PlayerId, times: p.Times}
		}
		return items, "top players by games", "top game player", nil
	}
}

func (h *topTab) buildKeyboard(activeTab string) service.InlineKeyboard {
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
		label := h.trans.Get(t.transKey, translator.Config{})
		if t.key == activeTab {
			label = fmt.Sprintf("▸ %s", label)
		}
		btn := service.KeyboardButton{
			Data:  TopTab.SetArgs(t.key).ToString(),
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

type topPlayerItem interface {
	PlayerID() int64
	Count() int
}

type simplePlayer struct {
	PID   int64
	times int
}

func (p simplePlayer) PlayerID() int64 { return p.PID }
func (p simplePlayer) Count() int      { return p.times }

type streakPlayer struct {
	PID                 int64
	ParticipationStreak int
	PeakStreak          int
}

func (p streakPlayer) PlayerID() int64 { return p.PID }
func (p streakPlayer) Count() int      { return p.ParticipationStreak }
