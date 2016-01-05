package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/geotrace/rest"
)

type ctxType byte // тип для сохранения данных в контексте запроса

var ctxToken ctxType = 1 // ключ для сохранения токенов

var (
	tokenType       = "sub"    // ключ токена, используемый для указания типа
	tokenTypeUser   = "user"   // токен с авторизацией пользователя
	tokenTypeDevice = "device" // токен авторизации устройства
)

// GetToken возвращает содержимое токена из контекста запроса.
func GetToken(c *rest.Context) map[string]interface{} {
	return c.Data(ctxToken).(map[string]interface{})
}

// Token проверяет токен, считывая его из заголовка. В случае неверного токена возвращает ошибку,
// что запрос не авторизован. Так же проверяет, что тип токена соответствует указанному в
// параметрах, в противном случае тоже будет ошибка. Сам токен сохраняется в контексте запроса.
func (s *Store) Token(h rest.Handler, allowSubs ...string) rest.Handler {
	return func(c *rest.Context) {
		token, err := s.tokens.ParseRequest(c.Request) // читаем токен из заголовка
		if err == jwt.ErrNoTokenInRequest {            // нет токена
			c.SetHeader("WWW-Authenticate", fmt.Sprintf("Bearer realm=%q", Realm))
			c.Status(http.StatusUnauthorized).Send(err)
			return
		}
		if err != nil { // токен не валиден
			c.Status(http.StatusForbidden).Send(err)
			return
		}
		if len(allowSubs) > 0 { // проверяем тип токена на допустимость
			tokenSub := token[tokenType]
			var allow bool
			for _, sub := range allowSubs {
				if tokenSub == sub {
					allow = true
					break
				}
			}
			if !allow { // токен не подходит под допустимый тип
				c.Status(http.StatusForbidden).Send("unauthorized token subject")
				return
			}
		}
		c.SetData(ctxToken, token) // сохраняем токен в контексте запроса
		h(c)                       // вызываем основной обработчик запроса
	}
}

// AccessTokenParamName описывает название параметра с токеном в запросе, если токен передается
// не в заголовке. Если в качестве значения задать пустую строку, то система перестанет
// поддерживать возможность передачи токена в виде параметра.
var AccessTokenParamName = "token"

const cryptoKeyLength = 1 << 8

// TokenEngine описывает класс для работы с токенами в формате JSON Web Token.
type TokenEngine struct {
	issuer    string        // название сервиса
	expire    time.Duration // время жизни ключа
	cryptoKey []byte        // ключ для подписи JWT
}

// NewTokenEngine инициализирует и возвращает класс для работы с токенами.
// Если ключ для подписи указан пустой, то формируется новый случайный ключ.
func NewTokenEngine(issuer string, expire time.Duration, cryptoKey []byte) *TokenEngine {
	if cryptoKey == nil {
		cryptoKey = make([]byte, cryptoKeyLength)
		if _, err := rand.Read(cryptoKey); err != nil {
			panic(err)
		}
	}
	return &TokenEngine{
		issuer:    issuer,
		expire:    expire,
		cryptoKey: cryptoKey,
	}
}

// Token формирует и возвращает токен в формате JWT.
func (e *TokenEngine) Token(items map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // генерируем новый токен
	for key, value := range items {          // добавляем в него наши данные
		token.Claims[key] = value
	}
	if e.issuer != "" { // добавляем информацию о сервисе
		token.Claims["iss"] = e.issuer
	}
	if e.expire != 0 { // время жизни токена
		token.Claims["exp"] = time.Now().Add(e.expire).Unix()
	}
	return token.SignedString(e.cryptoKey)
}

// verify является функцией для проверки целостности токена.
func (e *TokenEngine) verify(token *jwt.Token) (key interface{}, err error) {
	key = e.cryptoKey // ключ, используемый для подписи
	// проверяем метод вычисления сигнатуры и обязательные поля
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		err = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	} else if e.issuer != "" && token.Claims["iss"] != e.issuer {
		err = fmt.Errorf("unexpected Issuer: %v", token.Claims["iss"])
	} else if _, ok := token.Claims["exp"]; e.expire != 0 && !ok {
		err = errors.New("missing Expire")
	}
	return
}

// Parse разбирает токен, проверяет его валидность и возвращает данные из него.
func (e *TokenEngine) Parse(tokenString string) (data map[string]interface{}, err error) {
	token, err := jwt.Parse(tokenString, e.verify)
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}

// ParseRequest разбирает токен из HTTP-запроса. Токен может быть передан как в заголовке
// запроса авторизации, с типом авторизации "Bearer", так и в параметре или поле формы
// с имененен, определеннов в глобальной переменной AccessTokenParamName.
func (e *TokenEngine) ParseRequest(req *http.Request) (data map[string]interface{}, err error) {
	if ah := req.Header.Get("Authorization"); ah != "" {
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return e.Parse(ah[7:])
		}
	}
	if AccessTokenParamName != "" {
		if tokStr := req.FormValue(AccessTokenParamName); tokStr != "" {
			return e.Parse(tokStr)
		}
	}
	return nil, jwt.ErrNoTokenInRequest
}

// CryptoKey возвращает ключ, используемый для подписи, в виде строки (base64-encoded).
func (e *TokenEngine) CryptoKey() string {
	return base64.StdEncoding.EncodeToString(e.cryptoKey)
}
