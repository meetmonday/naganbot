package config

import "os"

const (
	TelegramBotToken = "TELEGRAM_BOT_TOKEN"
	DatabaseDsn      = "DATABASE_DSN"
	ForceUsername    = "FORCE_USERNAME"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}
