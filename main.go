package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/geotrace/rest"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	log         = log15.New()           // вывод логов
	retry       = 5                     // количество попыток подключения к сервисам
	delay       = time.Second           // время задержки между подключениями к сервисам в случае ошибки
	Realm       = "GeoTrace"            // используется в заголовке авторизации
	TokenIssuer = "com.xyzrd.geotracer" // автор токена
	TokenExpire = time.Minute * 30      // время жизни токена
)

func main() {
	// инициализируем параметры и окружение
	mongoURL := flag.String("mongodb", Env("MONGODB", "mongodb://localhost/geotrace"),
		"MongoDB connection `URL`")
	addr := flag.String("http", Env("SERVER", ":8080"), "HTTP server `address:port`")
	flag.Parse()

	// подключаемся к MongoDB и сервисам
	store, err := Connect(*mongoURL)
	if err != nil {
		log.Error("Connection error", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	// определяем обработчики URL
	var mux rest.ServeMux
	mux.Handles(rest.Paths{
		"user": {
			"GET":  Basic(store.UserLogin), // авторизация пользователя
			"POST": nil,                    // регистрация нового пользователя
		},
		"device": {
			"GET":  Basic(store.DeviceLogin), // авторизация устройства
			"POST": nil,                      // регистрация нового устройства
		},
		"users": {
			"GET": nil, // отдает список пользователей в группе
		},
		"users/:user-id": {
			"GET": nil, // отдает информацию о пользователе
		},
		"devices": {
			"GET": nil, // список устройств в группе
		},
		"devices/:device-id": {
			"GET":    nil, // информация об устройстве
			"PUT":    nil, // изменяет устройство
			"DELETE": nil, // удаляет
		},
		"devices/:device-id/events": {
			"GET":  nil, //
			"POST": nil, //
		},
		"devices/:device-id/events/:event-id": {
			"GET":    nil, //
			"PUT":    nil, //
			"DELETE": nil, //
		},
		"places": {
			"GET":  nil, //
			"POST": nil, //
		},
		"places/:place-id": {
			"GET":    nil, //
			"PUT":    nil, //
			"DELETE": nil, //
		},
	})
	mux.BasePath = "/api/v1/"
	server := http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("HTTP Server error", "err", err)
	}
}

// Env получает значение из окружения с заданным именем. Если значение не установлено, то
// возвращает значение, заданное по умолчанию.
func Env(envKey, defaultValue string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	return defaultValue
}
