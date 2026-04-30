package app

import (
	"context"
	"task-tracker-telegram-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	bot    *tgbotapi.BotAPI
	cfg    *config.Config
}

func NewApp(cfg *config.Config, logger *zap.Logger) (*App, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	logger.Info("Bot successfully initialized",
		zap.String("username", bot.Self.UserName),
	)
	return &App{
		logger: logger,
		bot:    bot,
		cfg:    cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := a.bot.GetUpdatesChan(u)

	a.logger.Info("bot started polling")

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("shutting down gracefully")
			return ctx.Err()
		case update, ok := <-updates:
			if !ok {
				a.logger.Info("updates channel closed")
				return nil
			}
			a.handleUpdate(update)
		}
	}
}

func (a *App) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("panic recovered in handleUpdate",
				zap.Any("recover", r),
				zap.Stack("stack"),
			)
		}
	}()

	if update.Message == nil {
		a.logger.Debug("update without message received")
		return
	}

	if update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я жалкое подобие кайтена, что ты можешь сделать:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить задание", "add_task"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Удалить задание", "delete_task"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Получить все задания", "get_all_tasks"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Получить конкретное задание", "get_task"),
			),
		)

		if _, err := a.bot.Send(msg); err != nil {
			a.logger.Error("failed to send message",
				zap.Error(err),
				zap.Int64("chat_id", update.Message.Chat.ID),
				zap.Int64("user_id", update.Message.From.ID),
				zap.String("message_type", "keyboard_response"),
			)
			return
		}
	}

	a.logger.Debug("message sent successfully",
		zap.Int64("chat_id", update.Message.Chat.ID),
	)
}
