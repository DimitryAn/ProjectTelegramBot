package tg

import (
	"bot/internal/lib/errWrap"
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
)

const (
	deleteAll      = true
	deleteSpecific = false
	singleLimit    = 1
	currentLimit   = 3
)

type Sender interface {
	SendMessage(chatID int, text string) error
}

type Operation interface {
	// Метод для сохранения новых заметок
	Save(ctx context.Context, text string, userName string) error

	// Метод для удаления заметок по полю text
	// При необходимости можно удалить все записи, для этого
	// необоходимо передать all = true
	Delete(ctx context.Context, userName string, text string, all bool) error

	// Извлечение заметок
	Extract(ctx context.Context, userName string, cnt int) ([]string, error)
}

type Processor struct {
	client Sender
	db     Operation
}

// Инициализация процессора
func NewProcessor(s Sender, db Operation) *Processor {
	return &Processor{
		client: s,
		db:     db,
	}
}

// Обработка команды от пользователя
func (p *Processor) MakeResponse(ctx context.Context, text string, chatID int, userName string) error {

	if text != "" && text[0] == '/' {
		text = strings.TrimSpace(text)
	}

	switch text {
	case "/start":
		err := p.client.SendMessage(chatID, StartCommand)
		if err != nil {
			return errWrap.Wrap("/start", err)
		}
	case "/delete":
		err := p.db.Delete(ctx, userName, text, deleteAll)
		if errors.Is(err, sql.ErrNoRows) {
			_ = p.client.SendMessage(chatID, EmptyPageMessage)
			return nil
		}
		if err != nil {
			return errWrap.Wrap("can't delete text (makeResponse)", err)
		}
		_ = p.client.SendMessage(chatID, DeleteCommand)
	case "/check":
		dates, err := p.db.Extract(ctx, userName, singleLimit)

		if err != nil {
			return errWrap.Wrap("can't check text (makeResponse)", err)
		}

		if len(dates) == 0 {
			_ = p.client.SendMessage(chatID, EmptyPageMessage)
			return nil
		}
		_ = p.client.SendMessage(chatID, dates[0])
		err = p.db.Delete(ctx, userName, dates[0], deleteSpecific)

		if err != nil {
			return errWrap.Wrap("can't delete page (makeResponse)", err)
		}
	case "/check3":
		dates, err := p.db.Extract(ctx, userName, currentLimit)
		if err != nil {
			return errWrap.Wrap("can't check3 (makeResponse)", err)
		}

		if len(dates) == 0 {
			_ = p.client.SendMessage(chatID, EmptyPageMessage)
			return nil
		}

		for _, data := range dates {
			_ = p.client.SendMessage(chatID, data)
			_ = p.db.Delete(ctx, userName, data, deleteSpecific)
		}
	case "/help":
		err := p.client.SendMessage(chatID, HelpCommand)
		if err != nil {
			log.Print(err)
		}
	default:
		err := p.db.Save(ctx, text, userName)
		if err != nil {
			return errWrap.Wrap("can't save text (makeResponse)", err)
		}
		_ = p.client.SendMessage(chatID, SaveMessage)
	}

	return nil
}
