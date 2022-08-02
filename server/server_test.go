package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	testingData := []struct {
		Title, Username, Password string
		ExpectedStatusCode        int
	}{
		{
			Title:              "Non-existing Username",
			Username:           "random",
			Password:           "random",
			ExpectedStatusCode: 400,
		},
		{
			Title:              "Existing Username",
			Username:           "admin",
			Password:           "admin",
			ExpectedStatusCode: 200,
		},
	}

	for _, test := range testingData {
		t.Run(test.Title, func(t *testing.T) {
			requestBody := strings.NewReader(fmt.Sprintf("username=%s&password=%s", test.Username, test.Password))
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/login", requestBody)
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()
			LoginHandler(recorder, request)

			result := recorder.Result()

			response, _ := io.ReadAll(result.Body)
			body := string(response)
			fmt.Println(body)

			require.Equal(t, test.ExpectedStatusCode, result.StatusCode)
		})
	}
}
