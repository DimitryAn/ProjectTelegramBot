package tg_test

import (
	"bot/internal/tg"
	"bot/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandStart(t *testing.T) {

	mockDb := mocks.NewMockOperation(t)
	mockSender := mocks.NewMockSender(t)
	mockSender.On("SendMessage", 123, tg.StartCommand).Return(nil)

	pr := tg.NewProcessor(mockSender, mockDb, context.Background())

	err := pr.MakeResponse("/start", 123, "u")

	assert.Nil(t, err)
	mockDb.AssertExpectations(t)
	mockSender.AssertExpectations(t)

}
