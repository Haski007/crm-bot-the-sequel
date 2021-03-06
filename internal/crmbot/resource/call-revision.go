package resource

import (
	"fmt"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"

	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) callRevision(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	message := fmt.Sprintf("%s *Ревизия* %s", emoji.MagnifyingGlass, emoji.MagnifyingGlass)

	answer := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		message,
		keyboards.Revision)
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
