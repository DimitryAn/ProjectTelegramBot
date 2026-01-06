package third_party

type Fetcher interface {
	// Метод для извлечения сообщений из чата телеграмма
	FetchMessage() ([]Message, error)
}

// Необходимые поля для обработки запроса пользователя
type Message struct {
	IsMessage bool
	ChatID    int
	Username  string
	Text      string
}
