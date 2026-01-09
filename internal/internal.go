package internal

// Необходимые поля для обработки запроса пользователя
type Message struct {
	IsMessage bool
	ChatID    int
	Username  string
	Text      string
}
