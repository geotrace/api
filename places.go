package main

import (
	"net/http"

	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

func (s *Store) PlacesList(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	places, err := (*model.Places)(s.db).List(token.Group)
	if err == model.ErrNotFound {
		return c.Send(rest.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return c.Send(places)
}

func (s *Store) PlaceGet(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	place, err := (*model.Places)(s.db).Get(token.Group, c.Param("place-id"))
	if err == model.ErrNotFound {
		return c.Send(rest.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return c.Send(place)
}

func (s *Store) PlaceAdd(c *rest.Context) error {
	place := new(model.Place)
	if err := c.Bind(place); err != nil {
		return err
	}
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	if err := (*model.Places)(s.db).Create(token.Group, place); err != nil {
		if err == model.ErrBadPlaceData {
			return c.Error(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.Status(http.StatusCreated).Send(rest.JSON{"id": place.ID})
}

func (s *Store) PlaceDelete(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	if err := (*model.Places)(s.db).Delete(token.Group, c.Param("place-id")); err != nil {
		if err == model.ErrNotFound {
			return c.Send(rest.ErrNotFound)
		}
		return err
	}
	return c.Send(nil)
}

func (s *Store) PlaceChange(c *rest.Context) error {
	place := new(model.Place)
	if err := c.Bind(place); err != nil {
		return err
	}
	place.ID = c.Param("place-id")
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	if err := (*model.Places)(s.db).Update(token.Group, place); err != nil {
		if err == model.ErrNotFound {
			return c.Send(rest.ErrNotFound)
		}
		if err == model.ErrBadPlaceData {
			return c.Error(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.Send(nil)
}
