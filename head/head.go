package head

import (
	"bot/internal"
	"log"
	"sync"
	"time"
)

type Processor interface {
	// Метод для обработки запроса пользователя и отправки ответа в чат
	MakeResponse(text string, chatID int, userName string) error
}

type Fetcher interface {
	// Метод для извлечения сообщений из чата телеграмма
	FetchMessage() ([]internal.Message, error)
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
func (t *Tools) Work() {

	log.Print("Start work!")

	for {

		messeges, err := t.fetcher.FetchMessage()
		if len(messeges) != 0 {
			log.Print("get new message")
		}

		if err != nil {
			log.Print(err)
			continue
		}

		var wg sync.WaitGroup

		for _, msg := range messeges {
			log.Printf("message - %s, from - %s", msg.Text, msg.Username)

			wg.Go(func() {
				err := t.processor.MakeResponse(msg.Text, msg.ChatID, msg.Username)
				if err != nil {
					log.Print(err)
				}
			})
		}

		wg.Wait()

		time.Sleep(1 * time.Second) //раз секунду получаем обновления

	}
}
