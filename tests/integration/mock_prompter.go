package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type basicAnswer struct {
	question string
	answer   string
	wasAsked bool
}
type selectAnswer struct {
	question string
	options  []string
	answer   string
	wasAsked bool
}
type MockPrompter struct {
	T             *testing.T
	basicAnswers  []*basicAnswer
	selectAnswers []*selectAnswer
}

func (m *MockPrompter) PromptBasic(label string, validator func(input string) error) *string {
	// Find a matching answer
	var answer *basicAnswer
	for _, ans := range m.basicAnswers {
		if ans.question == label && !ans.wasAsked {
			answer = ans
			break
		}
	}
	require.NotNilf(m.T, answer, "No matching answer for prompt '%s'", label)

	answer.wasAsked = true
	return &answer.answer
}

func (m *MockPrompter) PromptSelect(label string, items []string) *string {
	// Find a matching answer
	var answer *selectAnswer
	for _, ans := range m.selectAnswers {
		if ans.question == label && reflect.DeepEqual(items, ans.options) && !ans.wasAsked {
			answer = ans
			break
		}
	}
	require.NotNil(m.T, answer, "Failed to answer prompt '%s': no matching answer", label)

	answer.wasAsked = true
	return &answer.answer
}

func (m *MockPrompter) AnswerBasic(question string, answer string) {
	m.basicAnswers = append(m.basicAnswers, &basicAnswer{
		question, answer, false,
	})
}

func (m *MockPrompter) AnswerSelect(question string, expectedOptions []string, answer string) {
	m.selectAnswers = append(m.selectAnswers, &selectAnswer{
		question, expectedOptions, answer, false,
	})
}

func (m *MockPrompter) AssertAllPrompts() {
	for _, answer := range m.basicAnswers {
		assert.Truef(m.T, answer.wasAsked, "Prompt for '%s' was never asked", answer.question)
	}
}
