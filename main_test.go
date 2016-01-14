package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"
	"testing"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"

	"github.com/geotrace/model"
	"github.com/mdigger/jwt"
	"github.com/mdigger/rest"
)

var store *Store
var baseURL string
var usertoken []byte

const mongoURL = "mongodb://localhost/geotrace-test"

func TestMain(m *testing.M) {
	llog.SetHandler(log15.StreamHandler(os.Stdout, log15.JsonFormat()))
	// llog.Debug("init test")
	rest.Debug = true // возвращаем нормальное описание ошибок

	key := make([]byte, 1<<8) // создаем ключ для подписи токенов
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	tokenEngine := &TokenTemplate{ // инициализируем работу с токенами
		Template: jwt.Template{
			Issuer:  "com.xyzrd.geotrace",
			Expire:  time.Minute * 30,        // срок жизни
			Created: true,                    // добавлять время создания
			Signer:  jwt.NewSignerHS256(key), // подпись токена
		},
	}

	di, err := mgo.ParseURL(mongoURL)
	if err != nil {
		llog.Error("Bad MongoDB URL", "err", err)
		os.Exit(2)
	}
	session, err := mgo.DialWithInfo(di)
	if err != nil {
		llog.Error("Error MongoDB connection", "err", err)
		os.Exit(2)
	}
	// доступ к хранилищу данных
	store = &Store{model.InitDB(session, di.Database)}
	// инициализируем API
	mux := InitAPI(store, tokenEngine)
	// тестовый веб-сервер
	ts := httptest.NewServer(mux)
	// базовый путь для вызовов API
	baseURL = ts.URL + mux.BasePath
	// pretty.Println(mux)
	// запускаем тесты
	code := m.Run()
	// удаляем базу по окончании теста
	if err := session.DB(di.Database).DropDatabase(); err != nil {
		llog.Error("Error delete DB", "err", err)
	}
	store.Close() // закрываем соединение по окончании
	os.Exit(code) // возвращаем код окончания
}

// getUserToken возвращает токен пользователя
func getUserToken() (token []byte, err error) {
	if len(usertoken) > 0 {
		return usertoken, nil
	}
	if store == nil {
		return nil, errors.New("not connected to store")
	}
	// создаем тестового пользователя
	if err = (*model.Users)(store.db).Create(&model.User{
		Login:    "test",
		GroupID:  "test_group",
		Name:     "Test User",
		Password: model.NewPassword("test"),
	}); err != nil && !mgo.IsDup(err) {
		return nil, err
	}
	req, err := http.NewRequest("GET", baseURL+"user", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("test", "test")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	usertoken, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		var errJSON = make(rest.JSON)
		if err = json.Unmarshal(usertoken, &errJSON); err != nil {
			return nil, err
		}
		return nil, errors.New(errJSON["error"].(string))
	}
	return usertoken, nil
}

type TestRequest struct {
	Name   string
	Method string
	URL    string
	Data   interface{}
	Status int
}

var OutResponse bool = true

// request выводит в консоль запрос, делает его и потом выводит в консоль
// ответ на запрос.
func request(test TestRequest, token []byte) (*http.Response, error) {
	if test.Name != "" {
		fmt.Printf("#### %s\n", test.Name)
	}
	var body io.Reader = nil
	if test.Data != nil {
		jdata, err := json.MarshalIndent(test.Data, "", "  ")
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jdata)
	}
	req, err := http.NewRequest(test.Method, baseURL+test.URL, body)
	if err != nil {
		return nil, err
	}
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(token)))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if OutResponse {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("###### Request:\n```http\n%s\n```\n", dump)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if OutResponse {
		dump, err := httputil.DumpResponse(resp, false)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewReader(data))
		if len(data) > 0 {
			var buf = bytes.NewBuffer(dump)
			if ct := resp.Header.Get("Content-Type"); strings.Contains(ct, "application/json") {
				if err := json.Indent(buf, data, "", "  "); err != nil {
					return nil, err
				}
			} else {
				buf.Write(data)
			}
			dump = buf.Bytes()
		}

		fmt.Printf("###### Response:\n```http\n%s\n```\n", dump)
		fmt.Println(strings.Repeat("-", 40))
		if test.Status != resp.StatusCode {
			fmt.Printf("**ERROR:**: status %d != %d\n", test.Status, resp.StatusCode)
		}
		fmt.Println()
	}
	if test.Status != resp.StatusCode {
		return resp, fmt.Errorf("%q:\nstatus %d != %d", test.Name, test.Status, resp.StatusCode)
	}
	return resp, nil
}
