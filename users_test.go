package main

import (
	"encoding/json"
	"testing"

	"github.com/geotrace/model"
)

func TestUsers(t *testing.T) {
	token, err := getUserToken()
	if err != nil {
		t.Fatal(err)
	}

	test := TestRequest{
		"Получение списка пользователей без токена",
		"GET",
		"users",
		nil,
		401,
	}
	if _, err = request(test, nil); err != nil {
		t.Error(err)
	}
	resp, err := request(TestRequest{
		"Получение списка пользователей",
		"GET",
		"users",
		nil,
		200,
	}, token)
	if err != nil {
		t.Error(err)
	}
	var users []model.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Error(err)
	}
}
