package main

import (
	"bot/internal/clients/tgclient"
	"bot/internal/start"
	"bot/internal/storage/sqlite"
	"bot/internal/tg"
	"context"
	"flag"
	"log"
)

const (
	limit  = 100 //максимальное количество сообщений с телеграмма
	sqPath = "data/sqLite/storage.db"
)

func main() {

	//получение токенов
	host, token := mustFlags()

	//инициализация БД (sqlite)
	sqlDb, err := sqlite.New(sqPath)

	if err != nil {
		log.Fatal("can't create Database ", err)
	}
	err = sqlDb.Init(context.TODO())
	if err != nil {
		log.Fatal("can't create Database")
	}

	//инициализация клиента (сейчас - тг)
	client := tgclient.New(host, token)

	//инициализация фетчера (забирает сообщения из тг)
	fetcher := tg.NewFetcher(client, limit)

	//инициалищация процессора (работает с базой данных + обработка сообщений из тг)
	processor := tg.NewProcessor(client, sqlDb, context.TODO())

	// запуск цикла, управляет фетчером и процессором
	h := start.New(fetcher, processor)
	h.Work()
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
