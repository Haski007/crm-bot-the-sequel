package resource

import (
	"fmt"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/config"
	"strings"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *CrmBotService) commandCategoryRemove(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	if len(update.Message.Text) < 18 {
		bot.Errorf(chatID, "Wrong type of command!")
		return
	}
	categoryID := strings.ReplaceAll(update.Message.Text[len(update.Message.Text)-36:], "_", "-")

	if err := bot.ProductRepository.RemoveAllByCategoryID(categoryID); err != nil {
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[commandCategoryRemove] ProductRepository.RemoveAllByCategoryID | err: %s", err))
		return
	}

	if err := bot.CategoryRepository.RemoveByID(categoryID); err != nil {
		if err == repository.ErrDocDoesNotExist {
			bot.Errorf(chatID,
				"No category with such ID: \"%s\"", categoryID)
			return
		}
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		bot.ReportToTheCreator(
			fmt.Sprintf("[CategoryRepository.RemoveByID] categoryID: {%s} | err: %s", categoryID, err))
	}

	message := "Категория успешно удалена " + emoji.Basket
	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MainMenu
	answer.ParseMode = config.MarkdownParseMode
	bot.Bot.Send(answer)
}
