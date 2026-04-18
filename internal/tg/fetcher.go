package tg

import (
	"bot/internal/clients/tgclient"
	"bot/internal/lib/errWrap"
	"log"
)

const (
	UnknonwCommand = "Unknonw command"
)

type Reciver interface {
	Updates(limit int, offset int) ([]tgclient.Update, error)
}

type TgFetcher struct {
	offset int
	limit  int
	client Reciver
}

// Инициализация фетчера
func NewFetcher(client Reciver, limit int) *TgFetcher {
	return &TgFetcher{
		offset: 0,
		limit:  limit,
		client: client,
	}
}

// Сбор сообщений из чата телеграмма
func (tf *TgFetcher) FetchMessage() ([]Message, error) {
	updates, err := tf.client.Updates(tf.limit, tf.offset)
	if err != nil {
		return nil, errWrap.Wrap("can't get new updates (FetchMessage)", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	resultMessage := make([]Message, 0, len(updates))

	for _, upd := range updates {
		temp := parse(&upd)
		if temp.IsMessage {
			resultMessage = append(resultMessage, *temp)
		} else {
			log.Println(UnknonwCommand)
		}

	}
	tf.offset = updates[len(updates)-1].ID + 1
	return resultMessage, nil
}

// Обработка пришедшего сообщения
func parse(upd *tgclient.Update) *Message {
	if upd.Message == nil {
		return &Message{
			IsMessage: false,
		}
	}
	return &Message{
		IsMessage: true,
		ChatID:    upd.Message.Chat.ID,
		Username:  upd.Message.From.Username,
		Text:      upd.Message.Text,
	}

}
