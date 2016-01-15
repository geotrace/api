package main

import (
	"net/http"

	"github.com/geotrace/model"
	"github.com/mdigger/rest"
)

// PlacesList возвращает список мест, определенных для данной группы.
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

// PlaceGet возвращает описание конкретного места в данной группе.
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

// PlaceAdd добавляет описание нового места в группу.
func (s *Store) PlaceAdd(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	place := new(model.Place)
	if err := c.Bind(place); err != nil {
		return err
	}
	if err := (*model.Places)(s.db).Create(token.Group, place); err != nil {
		if err == model.ErrBadPlaceData {
			return c.Error(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.Status(http.StatusCreated).Send(rest.JSON{"id": place.ID})
}

// PlaceDelete удаляет описание места группы по его идентификатору.
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

// PlaceChange изменяет описание места в группе.
func (s *Store) PlaceChange(c *rest.Context) error {
	token := GetToken(c)
	if token == nil {
		return ErrBadToken
	}
	place := new(model.Place)
	if err := c.Bind(place); err != nil {
		return err
	}
	place.ID = c.Param("place-id")
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
