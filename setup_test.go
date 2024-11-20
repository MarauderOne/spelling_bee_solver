package main

import (
	"testing"
	"github.com/MarauderOne/spelling_bee_solver/dictionary_tools"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewAnswerList(t *testing.T) {
	answerList := createNewAnswerList()
	var resultCount int = answerList.Count()

	assert.NotEmpty(t, answerList)
	assert.IsType(t, dictionary_tools.MySimpleDict{}, *answerList)
	assert.Equal(t, 79586, resultCount)
	assert.Contains(t, answerList.Words, "AARDVARK")
}

func TestReviseAnswerList(t *testing.T) {
	answerList := createNewAnswerList()
	regexPattern := "^.*U.*$"

	answerList = reviseAnswerList(answerList, regexPattern)
	var resultCount int = answerList.Count()

	assert.NotEmpty(t, answerList)
	assert.IsType(t, dictionary_tools.MySimpleDict{}, *answerList)
	assert.Equal(t, 20891, resultCount)
	assert.Contains(t, answerList.Words, "ABACUS")
}

func TestNonAlpha(t *testing.T) {

	t.Run("Test letter", func(t *testing.T) {
		nonAlphaTest := nonAlpha("A")
		var exampleBool bool

		assert.Empty(t, nonAlphaTest)
		assert.IsType(t, exampleBool, nonAlphaTest)
		assert.Equal(t, false, nonAlphaTest)
	})
	t.Run("Test number", func(t *testing.T) {
		nonAlphaTest := nonAlpha("3")
		var exampleBool bool

		assert.NotEmpty(t, nonAlphaTest)
		assert.IsType(t, exampleBool, nonAlphaTest)
		assert.Equal(t, true, nonAlphaTest)
	})

	t.Run("Test symbol", func(t *testing.T) {
		nonAlphaTest := nonAlpha("&")
		var exampleBool bool

		assert.NotEmpty(t, nonAlphaTest)
		assert.IsType(t, exampleBool, nonAlphaTest)
		assert.Equal(t, true, nonAlphaTest)
	})

}
