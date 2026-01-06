package tg

import (
	"bot/clients/telegramClients"
	"bot/lib/errWrap"
	"bot/storage"
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

type Processer struct {
	client *telegramClients.Client
	db     storage.Operation
	ctx    context.Context
}

// Инициализация процессора
func NewProcesser(c *telegramClients.Client, db storage.Operation, ctx context.Context) *Processer {
	return &Processer{
		client: c,
		db:     db,
		ctx:    ctx,
	}
}

// Обработка команды от пользователя
func (p *Processer) MakeResponse(text string, chatID int, userName string) error {

	if text != "" && text[0] == '/' {
		strings.TrimSpace(text)
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
