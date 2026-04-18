package tg_test

import (
	"bot/internal/clients/tgclient"
	"bot/internal/tg"
	"bot/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	limit  = 100
	offset = 0
)

func TestFetchEmptyMessage(t *testing.T) {
	mockReciver := mocks.NewMockReciver(t)
	fetcher := tg.NewFetcher(mockReciver, limit)

	mockReciver.On("Updates", limit, offset).Return(nil, nil)

	mesg, err := fetcher.FetchMessage()

	assert.Nil(t, mesg)
	assert.Nil(t, err)

	mockReciver.AssertExpectations(t)
}

func TestFetchMessageWithError(t *testing.T) {
	mockReciver := mocks.NewMockReciver(t)
	fetcher := tg.NewFetcher(mockReciver, limit)

	mockReciver.On("Updates", limit, offset).Return(nil, assert.AnError)
	mesg, err := fetcher.FetchMessage()

	assert.Nil(t, mesg)
	assert.Error(t, err)

	mockReciver.AssertExpectations(t)
}

func TestFetchMessageSuccessfully(t *testing.T) {
	mockReciver := mocks.NewMockReciver(t)
	fetcher := tg.NewFetcher(mockReciver, limit)

	mockReciver.On("Updates", limit, offset).Return(func() []tgclient.Update {

		messg := tgclient.IncomingMessage{
			Text: "something",
			From: tgclient.From{
				Username: "user",
			},
			Chat: tgclient.Chat{
				ID: 1,
			},
		}

		return []tgclient.Update{{ID: 1, Message: &messg}}
	}(), nil)

	mesg, err := fetcher.FetchMessage()

	rtrn := []tg.Message{
		{
			IsMessage: true,
			ChatID:    1,
			Username:  "user",
			Text:      "something",
		},
	}

	assert.Equal(t, rtrn, mesg)
	assert.Nil(t, err)

	mockReciver.AssertExpectations(t)
}

func TestFetchNilMessage(t *testing.T) {
	mockReciver := mocks.NewMockReciver(t)
	fetcher := tg.NewFetcher(mockReciver, limit)

	mockReciver.On("Updates", limit, offset).Return([]tgclient.Update{{ID: 1, Message: nil}}, nil)

	mesg, err := fetcher.FetchMessage()

	assert.Empty(t, mesg)
	assert.Nil(t, err)

	mockReciver.AssertExpectations(t)
}
