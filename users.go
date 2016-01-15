package main

import (
	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

// UsersList возвращает список пользователей, которые входят в ту же группу.
func (s *Store) UsersList(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	users, err := (*model.Users)(s.db).List(token.Group)
	if err == model.ErrNotFound {
		return c.Send(rest.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return c.Send(users)
}
