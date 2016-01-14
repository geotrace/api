package main

import (
	"encoding/json"
	"testing"

	"github.com/geotrace/model"
)

func TestDevices(t *testing.T) {
	token, err := getUserToken()
	if err != nil {
		t.Fatal(err)
	}

	test := TestRequest{
		"Получение списка устройств в группе пользователя без токена",
		"GET",
		"devices",
		nil,
		401,
	}
	if _, err = request(test, nil); err != nil {
		t.Error(err)
	}
	resp, err := request(TestRequest{
		"Получение списка устройств в группе пользователя",
		"GET",
		"devices",
		nil,
		200,
	}, token)
	if err != nil {
		t.Error(err)
	}
	var devices []model.Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		t.Error(err)
	}
}
