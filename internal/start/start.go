package start

import (
	"bot/internal/tg"
	"context"
	"log"
	"sync"
	"time"
)

type Processor interface {
	// Метод для обработки запроса пользователя и отправки ответа в чат
	MakeResponse(ctx context.Context, text string, chatID int, userName string) error
}

type Fetcher interface {
	// Метод для извлечения сообщений из чата телеграмма
	FetchMessage() ([]tg.Message, error)
}

type Tools struct {
	fetcher   Fetcher
	processor Processor
}

// инициализация фетчера и процессора
func New(tf Fetcher, pr Processor) *Tools {
	return &Tools{
		fetcher:   tf,
		processor: pr,
	}
}

// Запуск бота
func (t *Tools) Work(ctx context.Context) {

	log.Println("Start work!")

	for {
		select {
		case <-ctx.Done():
			log.Println("stop work")
			return
		default:
		}

		messeges, err := t.fetcher.FetchMessage()
		if len(messeges) != 0 {
			log.Println("get new message")
		}

		if err != nil {
			log.Println(err)
			continue
		}

		var wg sync.WaitGroup

		for _, msg := range messeges {
			log.Printf("message - %s, from - %s \n", msg.Text, msg.Username)

			wg.Go(func() {
				err := t.processor.MakeResponse(ctx, msg.Text, msg.ChatID, msg.Username)
				if err != nil {
					log.Println(err)
				}
			})
		}

		wg.Wait()

		time.Sleep(1 * time.Second) //раз секунду получаем обновления

	}
}
