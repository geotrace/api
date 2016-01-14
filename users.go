package main

import (
	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

func (s *Store) UsersList(c *rest.Context) error {
	users, err := (*model.Users)(s.db).List(GetToken(c).Group)
	if err == model.ErrNotFound {
		return c.Send(rest.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return c.Send(users)
}
