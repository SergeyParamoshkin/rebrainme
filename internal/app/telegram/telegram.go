package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Telegram struct {
	Logger *zap.Logger
	Bot    *tgbotapi.BotAPI
	Config *Config
	Done   chan struct{}
}

func NewTelegram(logger *zap.Logger, config *Config) *Telegram {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil
	}

	bot.Debug = config.Debug

	logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	return &Telegram{
		Logger: logger,
		Bot:    bot,
		Config: config,
		Done:   make(chan struct{}, 1),
	}
}

func (t *Telegram) StopBot() {
	t.Bot.StopReceivingUpdates()

	<-t.Done
}

func (t *Telegram) StartBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)

	go func(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			if update.Message != nil { // If we got a message
				t.Logger.Info(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				_, err := t.Bot.Send(msg)
				if err != nil {
					t.Logger.Error("error sending message", zap.Error(err))
				}
			}
		}
		t.Done <- struct{}{}
	}(updates)
}
