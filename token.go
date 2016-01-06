package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/geotrace/jwt"
	"github.com/geotrace/rest"
	"gopkg.in/mgo.v2"
)

// TokenTemplate описывает шаблон для генерации токена.
type TokenTemplate struct {
	jwt.Template // шаблон токена
}

// Token описывает основное содержимое токена.
type Token struct {
	Type  string `json:"type"`
	Id    string `json:"id"`
	Group string `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ErrTokenNotFound описывает ошибку, что токен в запросе не найден.
var ErrTokenNotFound = errors.New("token not found")

// ParseRequest разбирает токен из HTTP-запроса.
func (t *TokenTemplate) ParseRequest(req *http.Request) (*Token, error) {
	var token = new(Token)
	if ah := req.Header.Get("Authorization"); ah != "" {
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			if err := t.Parse([]byte(ah[7:]), token); err != nil {
				return nil, err
			}
			return token, nil
		}
	}
	return nil, ErrTokenNotFound
}

// Token проверяет токен, считывая его из заголовка. В случае неверного токена возвращает ошибку,
// что запрос не авторизован. Так же проверяет, что тип токена соответствует указанному в
// параметрах, в противном случае тоже будет ошибка. Сам токен сохраняется в контексте запроса.
func (t *TokenTemplate) Token(h rest.Handler, allowSubs ...string) rest.Handler {
	return func(c *rest.Context) {
		token, err := t.ParseRequest(c.Request) // читаем токен из заголовка
		if err == ErrTokenNotFound {            // нет токена
			c.SetHeader("WWW-Authenticate", fmt.Sprintf("Bearer realm=%q", Realm))
			c.Status(http.StatusUnauthorized).Send(err)
			return
		}
		if err != nil { // токен не валиден
			c.Status(http.StatusForbidden).Send(err)
			return
		}
		if len(allowSubs) > 0 { // проверяем тип токена на допустимость
			var allow bool
			for _, sub := range allowSubs {
				if token.Type == sub {
					allow = true
					break
				}
			}
			if !allow { // токен не подходит под допустимый тип
				c.Status(http.StatusForbidden).Send("unauthorized token subject")
				return
			}
		}
		c.SetData(ctxType(1), token) // сохраняем токен в контексте запроса
		h(c)                         // вызываем основной обработчик запроса
	}
}

// Basic осуществляет HTTP Basic авторизацию и возвращает авторизационный токен.
func (t *TokenTemplate) Basic(auth func(login, password string) (*Token, error)) rest.Handler {
	return func(c *rest.Context) {
		login, password, ok := c.BasicAuth()
		if !ok {
			c.SetHeader("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", Realm))
			c.Status(http.StatusUnauthorized).Send(nil)
			return
		}
		token, err := auth(login, password)
		if err != nil {
			if err == mgo.ErrNotFound || err == ErrBadPassword {
				c.Status(http.StatusForbidden).Send(nil)
			} else {
				c.Error(err)
			}
			return
		}
		tokenData, err := t.Template.Token(token)
		if err != nil {
			c.Error(err)
			return
		}
		c.ContentType = "application/jwt"
		c.Send(tokenData)
		return
	}
}

type ctxType byte // тип для сохранения данных в контексте запроса

// GetToken возвращает содержимое токена из контекста запроса.
func GetToken(c *rest.Context) *Token {
	return c.Data(ctxType(1)).(*Token)
}
