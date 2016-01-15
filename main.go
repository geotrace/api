package main

import (
	"crypto/rand"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/mdigger/jwt"
	"github.com/mdigger/rest"
	_ "github.com/mdigger/rest/codex" // включаем поддержку разных форматов данных

	"gopkg.in/inconshreveable/log15.v2"
)

var (
	llog        = log15.New()           // вывод логов
	retry       = 5                     // количество попыток подключения к сервисам
	delay       = time.Second           // время задержки между подключениями к сервисам в случае ошибки
	Realm       = "GeoTrace"            // используется в заголовке авторизации
	TokenIssuer = "com.xyzrd.geotracer" // автор токена
	TokenExpire = time.Minute * 30      // время жизни токена
)

func InitAPI(store *Store, token *TokenTemplate) *rest.ServeMux {
	// var deviceMux ServeMux
	// deviceMux.Handles(rest.Paths{
	// 	"": {
	// 		// авторизация устройства
	// 		"GET": token.Basic(store.DeviceLogin),
	// 		// регистрация нового устройства
	// 		"POST": nil,
	// 	},
	// 	"events": {
	// 		"GET":  nil,
	// 		"POST": nil,
	// 	},
	// 	"places": {
	// 		// отдает список мест
	// 		"GET": rest.Handlers(token.GetToken("device"), store.PlacesList),
	// 	},
	// 	"users": {
	// 		// отдает список пользователей в группе
	// 		"GET": rest.Handlers(token.GetToken("device"), store.UsersList),
	// 	},
	// 	"token": {
	// 		"GET": nil,
	// 	},
	// })
	// deviceMux.BasePath = "/device"
	// определяем обработчики URL
	var mux rest.ServeMux
	mux.Handles(rest.Paths{
		"user": {
			// авторизация пользователя
			"GET": token.Basic(store.UserLogin),
			// регистрация нового пользователя
			"POST": nil,
		},
		"users": {
			// отдает список пользователей в группе
			"GET": token.Get(store.UsersList, "user"),
		},
		"devices": {
			// список устройств в группе
			"GET":  token.Get(store.DevicesList, "user"),
			"POST": nil,
		},
		"devices/:device-id": {
			// информация об устройстве
			"GET": nil,
			// изменяет устройство
			"PUT": nil,
			// удаляет устройство
			"DELETE": nil,
		},
		"devices/:device-id/events": {
			"GET":  nil,
			"POST": nil,
		},
		"devices/:device-id/events/:event-id": {
			"GET":    nil,
			"PUT":    nil,
			"DELETE": nil,
		},
		"places": {
			// отдает список мест
			"GET": token.Get(store.PlacesList, "user"),
			// создает новое место
			"POST": token.Get(store.PlaceAdd, "user"),
		},
		"places/:place-id": {
			// возвращает описание места
			"GET": token.Get(store.PlaceGet, "user"),
			// изменение информации о месте
			"PUT": token.Get(store.PlaceChange, "user"),
			// удаляет место из списка группы
			"DELETE": token.Get(store.PlaceDelete, "user"),
		},

		"device": {
			// авторизация устройства
			"GET": token.Basic(store.DeviceLogin),
			// регистрация нового устройства
			"POST": nil,
		},
		"device/events": {
			"GET":  nil,
			"POST": nil,
		},
		"device/places": {
			// отдает список мест
			"GET": token.Get(store.PlacesList, "device"),
		},
		"device/users": {
			// отдает список пользователей в группе
			"GET": token.Get(store.UsersList, "device"),
		},
		"device/token": {
			"GET": nil,
		},
	})
	mux.BasePath = "/api/v0/"
	// mux.Headers = map[string]string{
	// 	"Server": "GeoTrace Server",
	// }
	return &mux
}

func main() {
	// инициализируем параметры и окружение
	mongoURL := flag.String("mongodb", Env("MONGODB", "mongodb://localhost/geotrace"),
		"MongoDB connection `URL`")
	addr := flag.String("http", Env("SERVER", ":8080"), "HTTP server `address:port`")
	flag.Parse()

	key := make([]byte, 1<<8) // создаем ключ для подписи токенов
	if _, err := rand.Read(key); err != nil {
		llog.Error("Error generating token signer key", "err", err)
		os.Exit(1)
	}
	tokenEngine := &TokenTemplate{ // инициализируем работу с токенами
		Template: jwt.Template{
			Issuer:  "com.xyzrd.geotrace",
			Expire:  time.Minute * 30,        // срок жизни
			Created: true,                    // добавлять время создания
			Signer:  jwt.NewSignerHS256(key), // подпись токена
		},
	}

	store, err := Connect(*mongoURL) // подключаемся к MongoDB
	if err != nil {
		llog.Error("Connection error", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	mux := InitAPI(store, tokenEngine) // инициализируем API
	server := http.Server{             // инициализируем HTTP-сервер
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	if err := server.ListenAndServe(); err != nil {
		llog.Error("HTTP Server error", "err", err)
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
