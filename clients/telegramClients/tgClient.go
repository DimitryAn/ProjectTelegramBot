package telegramClients

import (
	"bot/lib/errWrap"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

// Инициализация клиента для телеграмма
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: makeBasePath(token),
		client:   http.Client{},
	}
}

func makeBasePath(token string) string {
	return "bot" + token
}

// Отправка сообщения от бота
func (c *Client) SendMessage(chatID int, text string) error {
	querry := url.Values{}
	querry.Add("chat_id", strconv.Itoa(chatID))
	querry.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, querry)

	if err != nil {
		return errWrap.Wrap("can't send message: ", err)
	}

	return nil
}

// Получение новых сообщений из телеграмма
func (c *Client) Updates(limit int, offset int) ([]Update, error) {
	querry := url.Values{}
	querry.Add("offset", strconv.Itoa(offset))
	querry.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, querry)

	if err != nil {
		return nil, errWrap.Wrap("can't get updates: ", err)
	}

	var result TgResponse

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil, errWrap.Wrap("can't parse json: ", err)
	}

	return result.Result, nil
}

// Обращение к api телеграмма с соответсвующим методом (method)
func (c *Client) doRequest(method string, querry url.Values) (data []byte, err error) {

	defer func() { err = errWrap.WrapIfErr("can't do request", err) }()
	uParam := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, uParam.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = querry.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
