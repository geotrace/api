package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/geotrace/model"
	"github.com/geotrace/rest"

	"gopkg.in/mgo.v2"
)

var ErrBadPassword = errors.New("bad password")

// Basic осуществляет HTTP Basic авторизацию
func Basic(auth func(login, password string) (string, error)) rest.Handler {
	return func(c *rest.Context) {
		login, password, ok := c.BasicAuth()
		if !ok {
			c.SetHeader("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", Realm))
			c.Status(http.StatusUnauthorized).Send(nil)
			return
		}
		token, err := auth(login, password)
		if err == nil {
			c.ContentType = "application/jwt"
			c.Send([]byte(token))
			return
		}
		if err == mgo.ErrNotFound || err == ErrBadPassword {
			c.Status(http.StatusForbidden).Send(nil)
		} else {
			c.Status(http.StatusInternalServerError).Send(err)
		}
	}
}

// UserLogin читает заголовок запроса с HTTP Basic авторизацией, проверяет пользователя
// по базе данных и отдает в ответ авторизационный ключ в формате JWT.
func (s *Store) UserLogin(login, password string) (string, error) {
	user, err := (*model.Users)(s.db).Login(login)
	if err != nil {
		return "", err
	}
	if !user.Password.Compare(password) {
		return "", ErrBadPassword
	}
	return s.tokens.Token(rest.JSON{
		tokenType: tokenTypeUser,
		"id":      user.Login,
		"group":   user.GroupID,
		"name":    user.Name,
	})
}

// DeviceLogin читает заголовок запроса с HTTP Basic авторизацией, проверяет устройство
// по базе данных и отдает в ответ авторизационный ключ в формате JWT.
func (s *Store) DeviceLogin(login, password string) (string, error) {
	device, err := (*model.Devices)(s.db).Login(login)
	if err != nil {
		return "", err
	}
	if !device.Password.Compare(password) {
		return "", ErrBadPassword
	}
	return s.tokens.Token(rest.JSON{
		tokenType: tokenTypeDevice,
		"id":      device.ID,
		"group":   device.GroupID,
		"name":    device.Name,
	})
}
