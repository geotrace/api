package main

import (
	"time"

	"github.com/geotrace/model"
	"gopkg.in/mgo.v2"
)

type Store struct {
	db *model.DB // хранилище
}

// Connect устанавливает соединение с MongoDB.
func Connect(url string) (*Store, error) {
	// устанавливаем соединение с MongoDB
	di, err := mgo.ParseURL(url)
	if err != nil {
		return nil, err
	}
	var session *mgo.Session
	// делаем несколько попыток, если сразу не получилось
	for i := 1; i <= retry; i++ {
		session, err = mgo.DialWithInfo(di)
		if err == nil {
			break
		}
		if i < retry {
			time.Sleep(time.Duration(i) * delay)
		} else {
			return nil, err // это была последняя попытка
		}
	}
	// возвращаем инициализированное хранилище
	return &Store{
		db: model.InitDB(session, di.Database),
	}, nil
}

// Close закрывает соединение с MongoDB.
func (s *Store) Close() {
	s.db.Close()
}
