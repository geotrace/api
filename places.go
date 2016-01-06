package main

import (
	"net/http"

	"github.com/geotrace/model"
	"github.com/geotrace/rest"
	"gopkg.in/mgo.v2"
)

func (s *Store) GetPlaces(c *rest.Context) {
	places, err := (*model.Places)(s.db).List(GetToken(c).Group)
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
