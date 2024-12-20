package http_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/stretchr/testify/assert"
)

const (
	URL = "http://localhost:3000/api/v1"
)

var client *http.Client = &http.Client{}

func TestInvalidCheckInRoute(t *testing.T) {
	payload := map[string]interface{}{
		"email": "invalid-email",
	}

	bytePayload, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", URL+"/emails/check-in", bytes.NewBuffer(bytePayload))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, 400)
}

type CheckInResponse struct {
	Url string `json:"url"`
}

func TestValidCheckInRoute(t *testing.T) {
	var result CheckInResponse

	payload := map[string]interface{}{
		"email": fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)),
	}

	bytePayload, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", URL+"/emails/check-in", bytes.NewBuffer(bytePayload))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 201)

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}

	assert.NotNil(t, res.Body)
}

func TestInvalidAuthorizationRoute(t *testing.T) {
	token := "invalid-token"

	req, err := http.NewRequest("GET", URL+"/emails/authorize?token"+token, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 400)
}

func TestValidAuthorizationRoute(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsLmRldW1weXRob25pc3RhQGdtYWlsLmNvbSIsImV4cCI6MTczNDUyNjE4Mn0.Sv_w7Aw5bM3uvJNLorn4yTr-yzL1Nr4u_y3dzSF0EPA"

	req, err := http.NewRequest("GET", URL+"/emails/authorize?token"+token, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
}

func TestInvalidDismissRoute(t *testing.T) {

	req, err := http.NewRequest("DELETE", URL+"/emails/dismiss/invalid@email.com", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 500)
}

func TestValidDismissRoute(t *testing.T) {
	req, err := http.NewRequest("PUT", URL+"/emails/dismiss/email.deumpythonista@gmail.com", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
}
