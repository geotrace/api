# REST API 

[![GoDoc](https://godoc.org/github.com/geotrace/api?status.svg)](https://godoc.org/github.com/geotrace/api)
[![Build Status](https://travis-ci.org/geotrace/api.svg)](https://travis-ci.org/geotrace/api)
[![Coverage Status](https://coveralls.io/repos/geotrace/api/badge.svg?branch=master&service=github)](https://coveralls.io/github/geotrace/api?branch=master)

Поддержка REST API для сервиса

- `/user` 
	+ [x] `GET` - авторизация пользователя и получение токена для работы с другими методами API
- `/device`
	+ [x] `GET` - авторизация устройства и получение токена для работы с другими методами API
- `/users`
	+ [x] `GET` - возвращает список пользователей
- `/devices`
	+ [x] `GET` - возвращает список устройств (только для пользователей)
- `/places`
	+ [x] `GET` - возвращает список мест для группы
	+ [x] `POST` - добавляет описание нового места (только для пользователей)
- `/places/{place-id}`
	+ [x] `GET` - возвращает информацию о месте
	+ [x] `PUT` - изменяет информацию о месте (только для пользователей)
	+ [x] `DELETE` - удаляет информацию о метсе (только для пользователей)