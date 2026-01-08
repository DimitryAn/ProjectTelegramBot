package head

import (
	"bot/internal"
	"log"
	"sync"
	"time"
)

type Tools struct {
	fetcher   internal.Fetcher
	processor internal.Processor
}

// инициализация фетчера и процессора
func New(tf internal.Fetcher, pr internal.Processor) *Tools {
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
