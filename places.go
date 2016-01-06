package main

import (
	"net/http"

	"github.com/geotrace/model"
	"github.com/geotrace/rest"
	"gopkg.in/mgo.v2"
)

func (s *Store) GetPlaces(c *rest.Context) {
	token := GetToken(c)
	if token == nil {
		c.Status(http.StatusForbidden).Send(nil)
		return
	}
	places, err := (*model.Places)(s.db).List(token.Group)
	if err != nil {
		if err == mgo.ErrNotFound {
			c.Status(http.StatusNotFound).Send(nil)
		} else {
			c.Error(err)
		}
		return
	}
	c.Send(places)
}
