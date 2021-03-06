package resource

import (
	"fmt"
	"strconv"

	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/model/keyboards"
	"github.com/Haski007/crm-bot-the-sequel/internal/crmbot/persistance/repository"
	"github.com/Haski007/crm-bot-the-sequel/pkg/emoji"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	getQuantitySetProduct Step = iota + 5
	getQuantitySetValue   Step = iota + 5
)

func (bot *CrmBotService) callQuantitySetHandler(update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	message := "Выберите категорию продукта:"

	var categories []string

	if err := bot.CategoryRepository.DistinctCategories(&categories); err != nil {
		bot.ReportToTheCreator(fmt.Sprintf("[callQuantitySet] DistinctCategories | err: %s", err))
		bot.Errorf(chatID,
			"Internal Server Error | write to @pdemian to get some help")
		return
	}

	if len(categories) == 0 {
		message := "Oops!"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)

		bot.Errorf(chatID, "Нет категорий!")
		return
	}

	OpsQueue[userID] = &Operation{
		Name: OperationType_QuantitySet,
		Step: getProductCategory,
		Data: nil,
	}

	answer := tgbotapi.NewMessage(chatID, message)
	answer.ReplyMarkup = keyboards.MarkupByArray(categories)
	bot.Bot.Send(answer)
}

func (bot *CrmBotService) hookQuantitySet(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	op := OpsQueue[userID]

	switch op.Step {
	case getProductCategory:
		categoryTitle := update.Message.Text

		var category model.Category
		if err := bot.CategoryRepository.FindByTitle(categoryTitle, &category); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Reply(chatID, fmt.Sprintf("Категории \"%s\" не существует! %s\n*Попробуй ещё раз*\n", categoryTitle, emoji.NoEntry))
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookQuantitySet] CategoryRepository.FindByTitle err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			delete(OpsQueue, userID)
			return
		}

		var products []string
		if err := bot.ProductRepository.FindTitlesByCategoryID(category.ID, &products); err != nil {
			if err == repository.ErrDocDoesNotExist {
				bot.Reply(chatID, fmt.Sprintf("Нет продуктов в этой категории, *попробуй другую* %s", emoji.NoEntry))
				return
			}
			bot.ReportToTheCreator(fmt.Sprintf("[hookQuantityAll] ProductRepository.FindTitlesByCategoryID | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			return
		}

		if len(products) == 0 {
			bot.Errorf(chatID,
				"В базе данных пока нет продуктов в этой категории")
			delete(OpsQueue, userID)
			return
		}

		OpsQueue[userID].Step++

		message := "Выберите продукт:"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MarkupByArray(products)
		bot.Bot.Send(answer)

	case getQuantitySetProduct:
		productTitle := update.Message.Text

		if !bot.ProductRepository.IsProductExists(productTitle) {
			bot.Reply(chatID, fmt.Sprintf("Продукта \"%s\" не существует! %s\n*Попробуй ещё раз*\n", productTitle, emoji.NoEntry))
			return
		}

		OpsQueue[userID].Data = productTitle
		OpsQueue[userID].Step++

		message := "Укажите новое количество продукта:"
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = tgbotapi.NewHideKeyboard(false)
		bot.Bot.Send(answer)
	case getQuantitySetValue:
		value, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			bot.Reply(chatID, fmt.Sprintf("\"%s\" - не натуральное число! %s\n*Попробуй ещё раз*\n", update.Message.Text, emoji.NoEntry))
			return
		}

		productTitle := OpsQueue[userID].Data.(string)

		if err := bot.ProductRepository.UpdateFieldByTitle(productTitle, "quantity", value); err != nil {
			bot.ReportToTheCreator(fmt.Sprintf("[hookQuantityAll] ProductRepository.UpdateFieldByTitle | err: %s", err))
			bot.Errorf(chatID,
				"Internal Server Error | write to @pdemian to get some help")
			return
		}

		message := "Количество продукта на складе успешно изменено " + emoji.Check
		answer := tgbotapi.NewMessage(chatID, message)
		answer.ReplyMarkup = keyboards.MainMenu
		bot.Bot.Send(answer)
	}
}
