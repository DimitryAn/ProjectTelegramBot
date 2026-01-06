package storage

import "context"

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
