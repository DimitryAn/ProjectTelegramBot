package head

import (
	"bot/third_party"
	"log"
	"sync"
	"time"
)

type Tools struct {
	fetcher third_party.Fetcher
	procces third_party.Processer
}

// инициализация фетчера и процессора
func New(tf third_party.Fetcher, pr third_party.Processer) *Tools {
	return &Tools{
		fetcher: tf,
		procces: pr,
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
				err := t.procces.MakeResponse(msg.Text, msg.ChatID, msg.Username)
				if err != nil {
					log.Print(err)
				}
			})
		}

		wg.Wait()

		time.Sleep(1 * time.Second) //раз секунду получаем обновления

	}
}
