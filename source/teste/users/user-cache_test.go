package users_test

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserCacheAuthz(t *testing.T) {
	client := resty.New()

	payload := map[string]string{
		"name":      "John",
		"last_name": "Doe",
		"cpf":       "123.456.789-12",
		"email":     "teste-cache@example.com",
		"password":  "senha1@AAs1$23",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/create")

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode())

	payload = map[string]string{
		"name":      "John",
		"last_name": "Doe",
		"cpf":       "123.456.789-13",
		"email":     "teste-cache-2@example.com",
		"password":  "senha1@AAs1$23",
	}

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/create")

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode())

	client = resty.New()

	email := "teste-cache@example.com"
	payload = map[string]string{
		"email":    email,
		"password": "senha1@AAs1$23",
	}

	respJWT, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/jwt")

	assert.NoError(t, err)
	assert.Equal(t, 200, respJWT.StatusCode())
	assert.Contains(t, respJWT.String(), "access_token")

	var jwtResponse map[string]interface{}
	err = json.Unmarshal(respJWT.Body(), &jwtResponse)
	assert.NoError(t, err)

	accessToken, ok := jwtResponse["access_token"].(string)
	assert.True(t, ok, "access_token not found in response")

	respMe, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		Get("http://localhost:8000/user/me")

	assert.NoError(t, err)
	assert.Equal(t, 200, respMe.StatusCode())

	var meResponse map[string]interface{}
	err = json.Unmarshal(respMe.Body(), &meResponse)
	assert.NoError(t, err)

	responseEmail, ok := meResponse["email"].(string)
	assert.True(t, ok, "Don't found email in /me response")
	assert.Equal(t, email, responseEmail, "The returned email don't match what was expected")

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		Delete("http://localhost:8000/user/del")

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())

	payload = map[string]string{
		"email":    "teste-cache-2@example.com",
		"password": "senha1@AAs1$23",
	}

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/jwt")

	json.Unmarshal(resp.Body(), &jwtResponse)
	accessToken, ok = jwtResponse["access_token"].(string)

	assert.NoError(t, err)
	assert.True(t, ok, "access_token not found in response")
	t.Log(jwtResponse)
	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+accessToken).
		Delete("http://localhost:8000/user/del")

	
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}
