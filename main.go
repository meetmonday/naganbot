package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/taranovegor/naganbot/container"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/handler/callback"
	"github.com/taranovegor/naganbot/handler/command"
	"github.com/taranovegor/naganbot/service"
	"github.com/taranovegor/naganbot/translator"
	"gorm.io/gorm"
)

var Version = "development"

func safeExecute(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()
	fn()
}

func main() {
	fmt.Println(fmt.Sprintf("Nagan bot! Version: %s", Version))

	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load dotenv file")
	}

	sc, err := container.Init()
	if err != nil {
		panic(err)
	}

	orm := sc.Get(container.ORM).(*gorm.DB)
	err = orm.AutoMigrate(
		&domain.Chat{},
		&domain.User{},
		&domain.Game{},
		&domain.Gunslinger{},
	)
	if err != nil {
		panic(err)
	}

	if orm.Migrator().HasColumn(&domain.Chat{}, "required_players") {
		orm.Exec("UPDATE chats SET required_players = 6 WHERE required_players IS NULL")
		orm.Migrator().AlterColumn(&domain.Chat{}, "Settings.RequiredPlayers")
	}

	if orm.Migrator().HasColumn(&domain.Game{}, "players_count") {
		orm.Exec("UPDATE games SET players_count = 6 WHERE players_count IS NULL")
		orm.Migrator().AlterColumn(&domain.Game{}, "PlayersCount")
	}

	if orm.Migrator().HasColumn(&domain.Game{}, "mode") {
		orm.Exec("UPDATE games SET mode = 'dynamic' WHERE mode IS NULL OR mode = ''")
		orm.Migrator().AlterColumn(&domain.Game{}, "Mode")
	}

	if orm.Migrator().HasColumn(&domain.Game{}, "status") {
		orm.Exec("UPDATE games SET status = 'lobby' WHERE status IS NULL OR status = ''")
		orm.Migrator().AlterColumn(&domain.Game{}, "Status")
	}

	if orm.Migrator().HasColumn(&domain.Chat{}, "settings_mode") {
		orm.Exec("UPDATE chats SET settings_mode = 'dynamic' WHERE settings_mode IS NULL OR settings_mode = ''")
		orm.Migrator().AlterColumn(&domain.Chat{}, "Settings.Mode")
	}

	if orm.Migrator().HasColumn(&domain.Chat{}, "last_hunger_post_at") {
		orm.Migrator().DropColumn(&domain.Chat{}, "last_hunger_post_at")
	}

	if !orm.Migrator().HasColumn(&domain.Chat{}, "hunger_stage") {
		orm.Migrator().AddColumn(&domain.Chat{}, "HungerStage")
	}

	trans := sc.Get(container.Translator).(*translator.Translator)
	chatRepository := sc.Get(container.RepositoryChat).(domain.ChatRepository)
	userRepository := sc.Get(container.RepositoryUser).(domain.UserRepository)
	bot := sc.Get(container.Bot).(*service.Bot)
	botApi := sc.Get(container.BotTelegram).(*tgbotapi.BotAPI)
	cmdRegistry := sc.Get(container.CommandRegistry).(*command.Registry)
	clbRegistry := sc.Get(container.CallbackRegistry).(*callback.Registry)
	scheduler := sc.Get(container.Scheduler).(*service.GameScheduler)

	log.Printf("authorized on account %s", botApi.Self.String())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		scheduler.Start(ctx)
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		chat := update.FromChat()
		if chat != nil {
			if chat.IsPrivate() || chat.IsChannel() {
				message := trans.Get("available only in chat", translator.Config{})
				bot.SendMessage(chat.ID, message)

				continue
			}

			domainChat := domain.NewChat(chat.ID, chat.Title, chat.UserName)
			if chatRepository.Exists(chat.ID) {
				if err := chatRepository.Update(domainChat); err != nil {
					log.Printf("failed to update chat %d: %v", chat.ID, err)
				}
			} else {
				if err := chatRepository.Store(domainChat); err != nil {
					log.Printf("failed to store chat %d: %v", chat.ID, err)
				}
			}
		}

		from := update.SentFrom()
		if from != nil {
			domainUser := domain.NewUser(from.ID, from.FirstName, from.LastName, from.UserName)
			if userRepository.Exists(from.ID) {
				if err := userRepository.Update(domainUser); err != nil {
					log.Printf("failed to update user %d: %v", from.ID, err)
				}
			} else {
				if err := userRepository.Store(domainUser); err != nil {
					log.Printf("failed to store user %d: %v", from.ID, err)
				}
			}
		}

		if update.Message != nil {
			msg := update.Message
			if !msg.IsCommand() {
				continue
			}
			name := msg.Command()
			cmd, err := cmdRegistry.Find(name)
			if err == nil {
				go safeExecute(func() { cmd.Execute(msg) })
			} else {
				log.Println(err.Error())
			}
		} else if update.CallbackQuery != nil {
			callbackQuery := update.CallbackQuery
			query := callback.Pattern(callbackQuery.Data)
			hdlr, err := clbRegistry.Find(query)
			if err == nil {
				go safeExecute(func() { hdlr.Execute(callbackQuery) })
			} else {
				log.Println(err.Error())
			}
		}
	}
}
