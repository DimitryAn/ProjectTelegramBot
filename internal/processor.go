package internal

type Processor interface {
	// Метод для обработки запроса пользователя и отправки ответа в чат
	MakeResponse(text string, chatID int, userName string) error
}
