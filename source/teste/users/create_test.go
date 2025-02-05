package users_test

import (
	"sync"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)


func TestCreateUser(t *testing.T) {
	client := resty.New()

	payload := map[string]string{
		"name":      "John",
		"last_name": "Doe",
		"cpf":       "123.456.789-11",
		"email":     "teste@example.com",
		"password":  "senha1@AAs1$23",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("http://localhost:8000/authn/create")

	assert.NoError(t, err)                   
	assert.Equal(t, 201, resp.StatusCode())  
}

func TestCreateUserConcurrent(t *testing.T) {
	client := resty.New()

	var wg sync.WaitGroup
	testCases := []struct {
		name     string
		lastName string
		cpf      string
		email    string
		password string
	}{
		{"John", "Doe", "123.456.789-11", "teste@example.com", "senha1@AAs1$23"},
		{"John", "Doe", "123.456.789-16", "teste1@example.com", "senha1@AAs1$23"},
		{"John", "Doe", "123.456.789-15", "teste2@example.com", "senha1@AAs1$23"},
		{"J", "Doe", "123.456.789-17", "teste3@example.com", "senha1@AAs1$23"},
		{"John", "Doe", "123.456.789-18", "teste4@example.com", "senha1@AAs1$23"},
		{"John", "Doe", "123.456.789-99", "teste5@example.com", "senha1@AAs1$23"},
	}

	for _, testCase := range testCases {
		wg.Add(1)
		go func(testCase struct {
			name     string
			lastName string
			cpf      string
			email    string
			password string
		}) {
			defer wg.Done()

			payload := map[string]string{
				"name":      testCase.name,
				"last_name": testCase.lastName,
				"cpf":       testCase.cpf,
				"email":     testCase.email,
				"password":  testCase.password,
			}
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(payload).
				Post("http://localhost:8000/authn/create")

			assert.NoError(t, err)
			assert.Equal(t, 201, resp.StatusCode())
		}(testCase)
	}
	wg.Wait()
}
