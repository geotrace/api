# REST API 

[![GoDoc](https://godoc.org/github.com/geotrace/api?status.svg)](https://godoc.org/github.com/geotrace/api)
[![Build Status](https://travis-ci.org/geotrace/api.svg)](https://travis-ci.org/geotrace/api)
[![Coverage Status](https://coveralls.io/repos/geotrace/api/badge.svg?branch=master&service=github)](https://coveralls.io/github/geotrace/api?branch=master)

Поддержка REST API для сервиса

### для пользователей

- `/login` 
	+ [ ] `GET` - авторизация пользователя и получение токена для работы с другими методами API
- `/users`
	+ [x] `GET` - возвращает список пользователей
- `/devices`
	+ [x] `GET` - возвращает список устройств
- `/devices/{device_id}`
	+ [ ] `GET` - возвращает информацию об устройстве
- `/devices/{device_id}/events`
	+ [ ] `GET` - возвращает список событий для данного устройства
- `/places`
	+ [x] `GET` - возвращает список мест для группы
	+ [x] `POST` - добавляет описание нового места
- `/places/{place_id}`
	+ [x] `GET` - возвращает информацию о месте
	+ [x] `PUT` - изменяет информацию о месте
	+ [x] `DELETE` - удаляет информацию о метсе

### для устройств

- `/device`
	+ [x] `GET` - авторизация устройства и получение токена для работы с другими методами API
	+ [ ] `PUT` - изменение информации об устройстве
	+ [ ] `POST` - регистрация нового устройства
- `/device/events`
	+ [ ] `GET` - возвращает список событий для данного устройства
	+ [ ] `POST` - публикует новое событие для данного устройства
- `/device/places`
	+ [ ] `GET` - возвращает список мест для группы
- `/device/users`
	+ [x] `GET` - возвращает список пользователей
- `/device/messages`
	+ [ ] `GET` - возвращает список сообщений для данного устройства
	+ [ ] `POST` - отправляет сообщение всем пользователям
- `/device/token`
	+ [ ] `GET` - генерирует и возвращает токен для присоединения устройства в группу
