package users_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoginJWT(t *testing.T) {
	client := resty.New()

	payload := map[string]string{
		"email":    "teste@example.com",
		"password": "senha1@AAs1$23",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/jwt")

	assert.NoError(t, err)                  
	assert.Equal(t, 200, resp.StatusCode()) 
	assert.Contains(t, resp.String(), "access_token")     
	assert.Contains(t, resp.String(), "eyJhbGciOiJIUzI1NiIs")  
}