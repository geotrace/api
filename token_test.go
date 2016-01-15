package main

import (
	"net/http"
	"testing"

	"github.com/mdigger/rest"
)

func TestTokens(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	c := &rest.Context{
		Request: req,
	}
	for _, f := range []func(c *rest.Context) error{
		store.DevicesList,
		store.PlacesList,
		store.PlaceGet,
		store.PlaceAdd,
		store.PlaceDelete,
		store.PlaceChange,
		store.UsersList,
	} {
		if err := f(c); err != ErrBadToken {
			t.Error(err)
		}
	}
}
