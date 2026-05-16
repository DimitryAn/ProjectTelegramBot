package tg_test

import (
	"bot/internal/lib/errWrap"
	"bot/internal/tg"
	"bot/mocks"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandStart(t *testing.T) {

	tests := []struct {
		chatID    int
		testName  string
		text      string
		userName  string
		mockSetup func(m *mocks.MockSender)
		expected  error
	}{
		{
			chatID:   123,
			testName: "successfully send message",
			text:     "/start",
			userName: "u",
			mockSetup: func(m *mocks.MockSender) {
				m.On("SendMessage", 123, tg.StartCommand).Return(nil)
			},
			expected: nil,
		},
		{
			chatID:   123,
			testName: "send message and got error",
			text:     "/start ",
			userName: "u",
			mockSetup: func(m *mocks.MockSender) {
				m.On("SendMessage", 123, tg.StartCommand).Return(errors.New("some db problem"))
			},
			expected: errWrap.Wrap("/start", errors.New("some db problem")),
		},
		{
			chatID:   123,
			testName: "got command with space",
			text:     "/start  ",
			userName: "u",
			mockSetup: func(m *mocks.MockSender) {
				m.On("SendMessage", 123, tg.StartCommand).Return(nil)
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			mockDb := mocks.NewMockOperation(t)
			mockSender := mocks.NewMockSender(t)
			test.mockSetup(mockSender)
			pr := tg.NewProcessor(mockSender, mockDb, context.Background())

			err := pr.MakeResponse(test.text, test.chatID, test.userName)

			if !assert.Equal(t, test.expected, err) {
				fmt.Println(err, test.expected)
				fmt.Printf("Error in test %s, expected - %s, got - %s", test.testName, test.expected, err.Error())
			}
			mockDb.AssertExpectations(t)
			mockSender.AssertExpectations(t)
		})
	}
}

// func TestSaveWithErr(t *testing.T) {
// 	mockDb := mocks.NewMockOperation(t)
// 	mockSender := mocks.NewMockSender(t)
// 	mockDb.On("Save", context.Background(), "my message", "u").Return(assert.AnError)

// 	pr := tg.NewProcessor(mockSender, mockDb, context.Background())

// 	err := pr.MakeResponse("my message", 123, "u")

// 	assert.Error(t, err)
// 	mockDb.AssertExpectations(t)
// 	mockSender.AssertExpectations(t)
// }
