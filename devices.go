package main

import (
	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

func (s *Store) DevicesList(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	devices, err := (*model.Devices)(s.db).List(token.Group)
	if err == model.ErrNotFound {
		return c.Send(rest.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return c.Send(devices)
}
