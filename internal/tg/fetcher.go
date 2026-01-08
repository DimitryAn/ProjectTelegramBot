package tg

import (
	"bot/clients/telegramClients"
	"bot/internal"
	"bot/lib/errWrap"
	"log"
)

const (
	UnknonwCommand = "Unknonw command"
)

type TgFetcher struct {
	client *telegramClients.Client
	offset int
	limit  int
}

// Инициализация фетчера
func NewFetcher(client *telegramClients.Client, limit int) *TgFetcher {
	return &TgFetcher{
		client: client,
		offset: 0,
		limit:  limit,
	}
}

// Сбор сообщений из чата телеграмма
func (tf *TgFetcher) FetchMessage() ([]internal.Message, error) {
	updates, err := tf.client.Updates(tf.limit, tf.offset)
	if err != nil {
		return nil, errWrap.Wrap("can't get new updates (FetchMessage)", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	resultMessage := make([]internal.Message, 0, len(updates))

	for _, upd := range updates {
		temp := parse(&upd)
		if temp.IsMessage {
			resultMessage = append(resultMessage, *temp)
		} else {
			log.Print(UnknonwCommand)
		}

	}
	tf.offset = updates[len(updates)-1].ID + 1
	return resultMessage, nil
}

// Обработка пришедшего сообщения
func parse(upd *telegramClients.Update) *internal.Message {
	if upd.Message == nil {
		return &internal.Message{
			IsMessage: false,
		}
	}
	return &internal.Message{
		IsMessage: true,
		ChatID:    upd.Message.Chat.ID,
		Username:  upd.Message.From.Username,
		Text:      upd.Message.Text,
	}

}
