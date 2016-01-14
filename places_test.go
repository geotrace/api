package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/geotrace/geo"
	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

func TestPlaces(t *testing.T) {
	token, err := getUserToken()
	if err != nil {
		t.Fatal(err)
	}

	tests := []TestRequest{
		{
			"Создание нового места",
			"POST",
			"places",
			rest.JSON{
				"name": "test_place",
				"circle": rest.JSON{
					"center": geo.Point{88, 55},
					"radius": 500,
				},
			},
			201,
		},
		{
			"Ошибка создания нового места",
			"POST",
			"places",
			rest.JSON{
				"name": "test_bad_place",
			},
			400,
		},
		{
			"Ошибка удаления несуществующего места",
			"DELETE",
			"places/bad_place",
			nil,
			404,
		},
	}

	for _, test := range tests {
		if _, err := request(test, token); err != nil {
			t.Error(err)
		}
	}

	test := TestRequest{
		"Получение списка мест без токена",
		"GET",
		"places",
		nil,
		401,
	}
	if _, err = request(test, nil); err != nil {
		t.Error(err)
	}
	resp, err := request(TestRequest{
		"Получение списка мест",
		"GET",
		"places",
		nil,
		200,
	}, token)
	if err != nil {
		t.Error(err)
	}
	var places []model.Place
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		t.Error(err)
	}
	for _, place := range places {
		test := TestRequest{
			fmt.Sprintf("Изменение места %q", place.String()),
			"PUT",
			fmt.Sprintf("places/%s", place.ID),
			rest.JSON{
				"name": "test_place_2",
				"circle": rest.JSON{
					"center": geo.Point{89, 89},
					"radius": 250,
				},
			},
			204,
		}
		if _, err = request(test, token); err != nil {
			t.Error(err)
		}
		test = TestRequest{
			fmt.Sprintf("Получение места %q", place.String()),
			"GET",
			fmt.Sprintf("places/%s", place.ID),
			nil,
			200,
		}
		if _, err = request(test, token); err != nil {
			t.Error(err)
		}
		test = TestRequest{
			fmt.Sprintf("Удаление места %q", place.String()),
			"DELETE",
			fmt.Sprintf("places/%s", place.ID),
			nil,
			204,
		}
		if _, err = request(test, token); err != nil {
			t.Error(err)
		}
	}
}
