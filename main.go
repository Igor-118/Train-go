package main

import (
	"context"
	"log"
	"os"

	"tgtest/clients/openweather"
	"tgtest/currency"
	"tgtest/facts"
	"tgtest/handler"
	"tgtest/memes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	//"github.com/jackc/pgx/v5"
	"tgtest/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Error ping db")
	}

	conn.QueryRow(context.Background(), "select id from users")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	owClient := openweather.New(os.Getenv("OPENWEATHERAPI_KEY"))

	userRepo := repo.NewRepo(conn)
	memKey := os.Getenv("MEM_KEY")
	mem := memes.NewMem(memKey)
	factKey := os.Getenv("FACT_KEY")
	fact := facts.NewFact(factKey)
	currencyKey := os.Getenv("CURRENCY_KEY")
	currency := currency.NewCurrency(currencyKey)
	botHandler := handler.NewHandler(
		bot,
		owClient,
		userRepo,
		mem,
		fact,
		currency,
	)

	botHandler.Start()
	//currency.GetCurrency()
}
