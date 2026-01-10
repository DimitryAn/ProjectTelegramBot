package tg

import (
	"bot/clients/telegramClients"
	"bot/lib/errWrap"
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
	client *telegramClients.Client
	db     Operation
	ctx    context.Context
}

// Инициализация процессора
func NewProcessor(c *telegramClients.Client, db Operation, ctx context.Context) *Processor {
	return &Processor{
		client: c,
		db:     db,
		ctx:    ctx,
	}
}

// Обработка команды от пользователя
func (p *Processor) MakeResponse(text string, chatID int, userName string) error {

	if text != "" && text[0] == '/' {
		text = strings.TrimSpace(text)
	}

	switch text {
	case "/start":
		err := p.client.SendMessage(chatID, startCommand)
		if err != nil {
			return errWrap.Wrap("/start", err)
		}
	case "/delete":
		err := p.db.Delete(p.ctx, userName, text, deleteAll)
		if errors.Is(err, sql.ErrNoRows) {
			_ = p.client.SendMessage(chatID, emptyPageMessage)
			return nil
		}
		if err != nil {
			return errWrap.Wrap("can't delete text (makeResponse)", err)
		}
		_ = p.client.SendMessage(chatID, deleteCommand)
	case "/check":
		dates, err := p.db.Extract(p.ctx, userName, singleLimit)

		if err != nil {
			return errWrap.Wrap("can't check text (makeResponse)", err)
		}

		if len(dates) == 0 {
			_ = p.client.SendMessage(chatID, emptyPageMessage)
			return nil
		}
		_ = p.client.SendMessage(chatID, dates[0])
		err = p.db.Delete(p.ctx, userName, dates[0], deleteSpecific)

		if err != nil {
			return errWrap.Wrap("can't delete page (makeResponse)", err)
		}
	case "/check3":
		dates, err := p.db.Extract(p.ctx, userName, currentLimit)
		if err != nil {
			return errWrap.Wrap("can't check3 (makeResponse)", err)
		}

		if len(dates) == 0 {
			_ = p.client.SendMessage(chatID, emptyPageMessage)
			return nil
		}

		for _, data := range dates {
			_ = p.client.SendMessage(chatID, data)
			_ = p.db.Delete(p.ctx, userName, data, deleteSpecific)
		}
	case "/help":
		err := p.client.SendMessage(chatID, helpCommand)
		if err != nil {
			log.Print(err)
		}
	default:
		err := p.db.Save(p.ctx, text, userName)
		if err != nil {
			return errWrap.Wrap("can't save text (makeResponse)", err)
		}
		_ = p.client.SendMessage(chatID, saveMessage)
	}

	return nil
}
