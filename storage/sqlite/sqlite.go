package sqlite

import (
	"bot/lib/errWrap"
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqClient struct {
	db *sql.DB
}

// Инициализация клиента sqlite
func New(path string) (*SqClient, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, errWrap.Wrap("can't open database: ", err)
	}

	if err := db.Ping(); err != nil {
		return nil, errWrap.Wrap("can't connet to database", err)
	}

	return &SqClient{db: db}, nil
}

func (sc *SqClient) Save(ctx context.Context, text string, userName string) error {

	if err := sc.db.Ping(); err != nil {
		return errWrap.Wrap("can't connet to database (save)", err)
	}

	q := `INSERT INTO Pages (page,user_name) VALUES (?,?)`

	ins, err := sc.db.ExecContext(ctx, q, text, userName)

	if cntRows, _ := ins.RowsAffected(); err != nil || cntRows == 0 {
		return errWrap.Wrap("can't save message: ", err)
	}

	return nil
}

func (sc *SqClient) Delete(ctx context.Context, userName string, page string, deleteAll bool) error {

	if err := sc.db.Ping(); err != nil {
		return errWrap.Wrap("can't connet to database", err)
	}

	var err error

	if deleteAll {
		err = sc.deleteAllPages(ctx, userName)
	} else {
		err = sc.deleteFirstPage(ctx, userName, page)
	}

	if err != nil {
		return err
	}

	return nil
}

func (sc *SqClient) Extract(ctx context.Context, userName string, limit int) ([]string, error) {

	if err := sc.db.Ping(); err != nil {
		return nil, errWrap.Wrap("can't connet to database (extract)", err)
	}

	q := `SELECT page FROM Pages WHERE user_name = ? LIMIT ?`

	rows, err := sc.db.QueryContext(ctx, q, userName, limit)

	defer func() { _ = rows.Close() }()

	if err != nil {
		return nil, errWrap.Wrap("can't do select from DB", err)
	}

	var result []string
	for rows.Next() {
		var page string
		err = rows.Scan(&page)
		if err != nil {
			return nil, errWrap.Wrap("can't scan page", err)
		}

		result = append(result, page)
	}

	return result, nil
}

// Создание базы данных
func (sc *SqClient) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS Pages (page TEXT,user_name TEXT)`
	_, err := sc.db.ExecContext(ctx, q)

	if err != nil {
		return errWrap.Wrap("can't up database ", err)
	}
	return nil
}

// Удаление всех записей по нику пользователя
func (sc *SqClient) deleteAllPages(ctx context.Context, userName string) error {

	q := `DELETE FROM Pages WHERE user_name=?`

	res, err := sc.db.ExecContext(ctx, q, userName)

	if err != nil {
		return errWrap.Wrap("can't delete pages: ", err)
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Удаление записи, которая была занесена самой первой
func (sc *SqClient) deleteFirstPage(ctx context.Context, userName string, page string) error {

	q := `DELETE FROM Pages WHERE user_name=? AND page=?`

	res, err := sc.db.ExecContext(ctx, q, userName, page)

	if err != nil {
		return errWrap.Wrap("can't delete pages: ", err)
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}
