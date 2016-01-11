package main

import (
	"github.com/geotrace/model"
	"github.com/geotrace/rest"
)

func (s *Store) GetPlaces(c *rest.Context) error {
	places, err := (*model.Places)(s.db).List(GetToken(c).Group)
	if err != nil {
		return err
	}
	return c.Send(places)
}
