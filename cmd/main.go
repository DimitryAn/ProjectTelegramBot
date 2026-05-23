package main

import (
	"bot/internal/clients/tgclient"
	"bot/internal/start"
	"bot/internal/storage/sqlite"
	"bot/internal/tg"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const (
	limit  = 100 //максимальное количество сообщений с телеграмма
	sqPath = "data/sqLite/storage.db"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	//получение токенов
	host, token := mustFlags()

	//инициализация БД (sqlite)
	sqlDb, err := sqlite.New(sqPath)

	if err != nil {
		log.Fatal("can't create Database ", err)
	}
	err = sqlDb.Init(ctx)
	if err != nil {
		log.Fatal("can't create Database")
	}

	//инициализация клиента (сейчас - тг)
	client := tgclient.New(host, token)

	//инициализация фетчера (забирает сообщения из тг)
	fetcher := tg.NewFetcher(client, limit)

	//инициалищация процессора (работает с базой данных + обработка сообщений из тг)
	processor := tg.NewProcessor(client, sqlDb)

	// запуск цикла, управляет фетчером и процессором
	h := start.New(fetcher, processor)
	go h.Work(ctx)

	<-ctx.Done()
	time.Sleep(5 * time.Second)
	log.Print("shut down program")
}

// Обрабатывает флаги при запуске программы
// Необходимо передать токен от телеграмм бота и хост откуда брать новые сообщения
// хост телеграмма - 'api.telegram.org'
func mustFlags() (string, string) {
	token := flag.String("tgToken", "", "token needed for start bot")
	host := flag.String("host", "", "host for bot")

	flag.Parse()

	if token == nil || *token == "" {
		log.Fatal("Empty token!")
	}

	if host == nil || *host == "" {
		log.Fatal("Empty host!")
	}

	return *host, *token

}
