package handler

import (
	"fmt"
	"log"
	"math"

	//"strings"
	"context"
	"strconv"
	"strings"
	"tgtest/clients/openweather"
	"tgtest/currency"
	"tgtest/facts"
	"tgtest/memes"
	"tgtest/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gt "gopkg.gilang.dev/translator/v2"
)

type userRepository interface {
	GetUserCity(ctx context.Context, userID int64) (string, error)
	CreateUser(ctx context.Context, userID int64, username string) error
	UpdateCity(ctx context.Context, userID int64, city string) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	IsUserExistsByUsername(ctx context.Context, username string) (bool, error)
	IsUserExistsByID(ctx context.Context, ID int64) (bool, error)
	//GetCurrencyValue(ctx context.Context, ID int64, amount string) error
}

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	userRepo userRepository
	mem      memes.MemClient
	fact     facts.FactCliet
	currency currency.CurrencyClient
}

func NewHandler(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, userRepo userRepository, mem memes.MemClient, fact facts.FactCliet, currency currency.CurrencyClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		userRepo: userRepo,
		mem:      mem,
		fact:     fact,
		currency: currency,
	}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.handleUpdate(update)
	}
}

func (h *Handler) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	ctx := context.Background()

	if update.Message.IsCommand() {
		err := h.ensureUser(ctx, update)
		if err != nil {
			log.Println("error ensureUser: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		switch update.Message.Command() {
		case "currency":
			h.handleCurrency(ctx, update)
		case "fact":
			h.handleFact(ctx, update)
		case "mem":
			h.handleMem(update)
		case "time":
			h.handleTime(update)
		case "start":
			h.handleStart(update)
		case "city":
			h.handleSetCity(ctx, update)
			return
		case "weather":
			h.handleSendWeather(ctx, update)
			return
		default:
			h.handleUnknownCommand(update)
			return
		}
	} else {
		fmt.Printf("Не команда %s\n", update.Message.Text)
	}
}

func (h *Handler) handleCurrency(ctx context.Context, update tgbotapi.Update) {
	//var nazvanie, znachenie string
	input := update.Message.CommandArguments()
	fmt.Printf("РАЗОБРАТЬ НАДО БЫ %s\n", input)
	cleaned := strings.TrimSpace(input)
	split := strings.Split(cleaned, " ")
	//fmt.Println(split, "split")

	if len(split) != 3 {
		log.Println("Сообщение криво написано")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка неверно написана комманда, надо например вот так: 100 USD RUB")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	chtoPerevodim, err := currency.ParseCurrencyName(split[1])
	//fmt.Println(chtoPerevodim, "что переводим")
	//fmt.Println(split[2], "Во что переводим")
	if err != nil {
		log.Println("Ошибка в получении функции handleCurrancy: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка с названием валюты")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	//fmt.Println(chtoPerevodim, "что переводим")

	newCurrency, err := h.currency.GetCurrency(chtoPerevodim)
	//fmt.Println(newCurrency.Rates["RUB"], "рубль")
	if err != nil {
		log.Println("Ошибка в получении функции GetCurrency: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка с функцией GetCurrency")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	//fmt.Println("Введи название валюты")
	//nazvanie := tgbotapi.NewUpdate(0)
	//nazvanie.Timeout = 60
	// updates := h.bot.GetUpdatesChan(nazvanie)
	// var nazvanieDeneg string
	// for update := range updates {
	// 	if update.Message != nil && update.Message.Text != "" {
	// 		nazvanieDeneg = update.Message.Text

	// 		// Выводим текст в консоль сервера
	// 		//log.Printf("[%s] написал: %s", nazvanieDeneg)

	// 		// Пример реакции: отправляем этот же текст обратно
	// 		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели: "+userText)
	// 		// if _, err := h.bot.Send(msg); err != nil {
	// 		// 	log.Println("Ошибка отправки сообщения:", err)
	// 		// }
	// 	}
	// }

	voChtoPerevodim, err := currency.ParseCurrencyName(split[2])
	//fmt.Println(chtoPerevodim, " Во что переводим")
	if err != nil {
		log.Println("Ошибка в получении функции handleCurrancy: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Название валюты должно быть длинною в 3 символа")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	// amount := tgbotapi.NewUpdate(0)
	// //amount.Timeout = 60
	// updatesDeneg := h.bot.GetUpdatesChan(amount)
	// var kolichestvoDeneg string
	// for update := range updatesDeneg {
	// 	if update.Message != nil && update.Message.Text != "" {
	// 		kolichestvoDeneg = update.Message.Text

	// 		// Выводим текст в консоль сервера
	// 		//log.Printf("[%s] написал: %s", kolichestvoDeneg)

	// 		// Пример реакции: отправляем этот же текст обратно
	// 		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ввели: "+userText)
	// 		// if _, err := h.bot.Send(msg); err != nil {
	// 		// 	log.Println("Ошибка отправки сообщения:", err)
	// 		// }
	// 	}
	// }
	kolichestvoDeneg := split[0]
	curAmount, err := currency.ParseNumber(kolichestvoDeneg)
	if err != nil {
		log.Println("Ошибка в получении функции handleCurrancy: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка с названием валюты")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	var total int
	for valyuta, znachenie := range newCurrency.Rates {
		if valyuta == voChtoPerevodim {
			g, err := strconv.Atoi(curAmount)
			if err != nil {
				fmt.Errorf("не смог перевести текст в число")
			}
			total = int(znachenie) * g
			break
		} else {
			log.Println("Ошибка в получении функции handleCurrancy: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Нет такой валюты: %s\nВот доступные валюты для перевода: EUR, RUB, ADA, USD", voChtoPerevodim))
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s %s равны %d %s", kolichestvoDeneg, chtoPerevodim, total, voChtoPerevodim))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

// func(h *Handler) ESCHEREERROR(forLog, forMsg string, err error, update tgbotapi.Update) {
// 		log.Println("Ошибка в получении функции handleCurrancy: ", err)
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка с названием валюты")
// 		msg.ReplyToMessageID = update.Message.MessageID
// 		h.bot.Send(msg)
// 		return
// }

func (h *Handler) handleFact(ctx context.Context, update tgbotapi.Update) {
	newFact, err := h.fact.GetFacts()
	if err != nil {
		log.Println("Ошибка в получении факта в функции handleFuct: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	result, _ := gt.Translate(ctx, newFact.Fact, "ru")
	translate := result.Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вот факт: %s ", translate))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleMem(update tgbotapi.Update) {
	keyword := update.Message.CommandArguments()
	mem, err := h.mem.Mem(keyword)
	if err != nil {
		log.Println("error Ошибка с функцией Mem: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вот мем: %s ", mem))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleTime(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Текущее время: %v",
		time.Now().Format("15:04:05")))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
	return
}

func (h *Handler) handleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("что бы установить нужный город, напиши команду /city\nЧтобы узнать погоду, напиши команду /weather\n"))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
	return
}

func (h *Handler) handleSetCity(ctx context.Context, update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	err := h.userRepo.UpdateCity(ctx, update.Message.From.ID, city)
	if err != nil {
		log.Println("error userRepo.UpdateCity: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город %s сохранен", city))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleSendWeather(ctx context.Context, update tgbotapi.Update) {
	city, err := h.userRepo.GetUserCity(ctx, update.Message.From.ID)
	if err != nil {
		log.Println("error userRepo.GetUserCity: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	if city == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала установите город с помощью команды /city\nВот так: ' /city Москва'")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	coordinates, err := h.owClient.Coordinates(city)
	if err != nil {
		log.Printf("error owClient.Coordinates: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить координаты")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Printf("error owClient.Weather: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить погоду в этой местности")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Температура в %s: %d°C", city, int(math.Round(weather.Temp))),
	)
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *Handler) handleUnknownCommand(update tgbotapi.Update) {
	log.Printf("Unknown command[%s] %s", update.Message.From.UserName, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такая команда недоступна")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
	return
}

func (h *Handler) ensureUser(ctx context.Context, update tgbotapi.Update) error {
	userID := update.Message.From.ID
	username := update.Message.From.UserName

	exists, err := h.userRepo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error userRepo.IsUserExist: %w", err)
	}

	if !exists {
		err := h.userRepo.CreateUser(ctx, userID, username)
		if err != nil {
			return fmt.Errorf("error userRepo.CreateUser: %w", err)
		}
	}
	return nil
}
