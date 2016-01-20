package main

import (
	"errors"

	"github.com/geotrace/model"
)

var ErrBadPassword = errors.New("bad password")

// UserLogin читает заголовок запроса с HTTP Basic авторизацией, проверяет
// пользователя по базе данных и отдает в ответ авторизационный ключ в формате
// JWT.
func (s *Store) UserLogin(login, password string) (*Token, error) {
	user, err := (*model.Users)(s.db).Login(login)
	if err != nil {
		return nil, err
	}
	if !user.Password.Compare(password) {
		return nil, ErrBadPassword
	}
	return &Token{
		Type:  "user",
		Id:    user.Login,
		Group: user.GroupID,
		Name:  user.Name,
	}, nil
}

// DeviceLogin читает заголовок запроса с HTTP Basic авторизацией, проверяет
// устройство по базе данных и отдает в ответ авторизационный ключ в формате
// JWT.
func (s *Store) DeviceLogin(login, password string) (*Token, error) {
	device, err := (*model.Devices)(s.db).Login(login)
	if err != nil {
		return nil, err
	}
	if !device.Password.Compare(password) {
		return nil, ErrBadPassword
	}
	return &Token{
		Type:  "device",
		Id:    device.ID,
		Group: device.GroupID,
		Name:  device.Name,
	}, nil
}
