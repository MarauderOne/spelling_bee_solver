package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSolveSpellingBee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	webServer := gin.Default()
	webServer.POST("/letters", parseLetters)

	t.Run("Test valid input", func(t *testing.T) {
		//Define a valid grid input
		gridData := []CellData{
			{Character: "U", Color: "yellow"},
			{Character: "D", Color: "grey"},
			{Character: "T", Color: "grey"},
			{Character: "V", Color: "grey"},
			{Character: "E", Color: "grey"},
			{Character: "I", Color: "grey"},
			{Character: "N", Color: "grey"},
		}

		jsonData, _ := json.Marshal(gridData)
		req, _ := http.NewRequest(http.MethodPost, "/letters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		httpResponse := httptest.NewRecorder()
		webServer.ServeHTTP(httpResponse, req)

		assert.Equal(t, http.StatusOK, httpResponse.Code)

		var jsonResponse map[string]interface{}
		err := json.Unmarshal(httpResponse.Body.Bytes(), &jsonResponse)
		assert.NoError(t, err)

		assert.NotEmpty(t, jsonResponse["result"])
		assert.NotEmpty(t, jsonResponse["resultCount"])
		assert.Contains(t, jsonResponse["result"], "UNINTENDED")
		assert.Equal(t, float64(39), jsonResponse["resultCount"])
	})

	t.Run("Test minimal input", func(t *testing.T) {
		//Define a valid grid input
		gridData := []CellData{
			{Character: "O", Color: "yellow"},
			{Character: "N", Color: "grey"},
		}

		jsonData, _ := json.Marshal(gridData)
		req, _ := http.NewRequest(http.MethodPost, "/letters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		httpResponse := httptest.NewRecorder()
		webServer.ServeHTTP(httpResponse, req)

		assert.Equal(t, http.StatusOK, httpResponse.Code)

		var jsonResponse map[string]interface{}
		err := json.Unmarshal(httpResponse.Body.Bytes(), &jsonResponse)
		assert.NoError(t, err)

		assert.NotEmpty(t, jsonResponse["result"])
		assert.NotEmpty(t, jsonResponse["resultCount"])
		assert.Equal(t, "NOON", jsonResponse["result"])
		assert.Equal(t, float64(1), jsonResponse["resultCount"])
	})

	t.Run("Test invalid character input", func(t *testing.T) {
		//Define an invalid grid input (invalid color)
		gridData := []CellData{
			{Character: "&", Color: "grey"},
		}

		jsonData, _ := json.Marshal(gridData)
		req, _ := http.NewRequest(http.MethodPost, "/letters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		httpResponse := httptest.NewRecorder()
		webServer.ServeHTTP(httpResponse, req)

		assert.Equal(t, http.StatusBadRequest, httpResponse.Code)

		var jsonResponse map[string]interface{}
		err := json.Unmarshal(httpResponse.Body.Bytes(), &jsonResponse)
		assert.NoError(t, err)

		assert.Contains(t, jsonResponse["error"], "Invalid character: &")
	})

	t.Run("Test invalid color input", func(t *testing.T) {
		//Define an invalid grid input (invalid color)
		gridData := []CellData{
			{Character: "P", Color: "purple"},
			{Character: "", Color: ""},
			{Character: "", Color: ""},
			{Character: "", Color: ""},
			{Character: "", Color: ""},
		}

		jsonData, _ := json.Marshal(gridData)
		req, _ := http.NewRequest(http.MethodPost, "/letters", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		webServer.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var jsonResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
		assert.NoError(t, err)

		assert.Contains(t, jsonResponse["error"], "Invalid color: purple")
	})
}
