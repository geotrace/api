## Создание нового места

###### Request:
```http
POST /api/v1/places HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I
Content-Type: application/json

{
    "circle": {
        "center": [
            88.000000,
            55.000000
        ],
        "radius": 500
    },
    "name": "test_place"
}
```
###### Response:
```http
HTTP/1.1 201 Created
Content-Length: 29
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT

{
	"id": "VpOFLDRe2bzqO6Ha"
}
```
----------------------------------------

## Ошибка создания нового места

###### Request:
```http
POST /api/v1/places HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I
Content-Type: application/json

{
    "name": "test_bad_place"
}
```
###### Response:
```http
HTTP/1.1 400 Bad Request
Content-Length: 76
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT

{
	"code": 400,
	"message": "bad place data: cyrcle or polygon is require"
}
```
----------------------------------------

## Ошибка удаления несуществующего места

###### Request:
```http
DELETE /api/v1/places/bad_place HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I


```
###### Response:
```http
HTTP/1.1 404 Not Found
Content-Length: 41
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT

{
	"code": 404,
	"message": "Not Found"
}
```
----------------------------------------

## Получение списка мест без токена

###### Request:
```http
GET /api/v1/places HTTP/1.1
Host: 127.0.0.1:62441


```
###### Response:
```http
HTTP/1.1 401 Unauthorized
Content-Length: 44
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT
Www-Authenticate: Bearer realm="GeoTrace"

{
	"code": 401,
	"message": "Unauthorized"
}
```
----------------------------------------

## Получение списка мест

###### Request:
```http
GET /api/v1/places HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I


```
###### Response:
```http
HTTP/1.1 200 OK
Content-Length: 146
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT

[
	{
		"id": "VpOFLDRe2bzqO6Ha",
		"name": "test_place",
		"circle": {
			"center": [
				88.000000,
				55.000000
			],
			"radius": 500
		}
	}
]
```
----------------------------------------

## Изменение места "test_place"

###### Request:
```http
PUT /api/v1/places/VpOFLDRe2bzqO6Ha HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I
Content-Type: application/json

{
    "circle": {
        "center": [
            89.000000,
            89.000000
        ],
        "radius": 250
    },
    "name": "test_place_2"
}
```
###### Response:
```http
HTTP/1.1 204 No Content
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT


```
----------------------------------------

## Получение места "test_place"

###### Request:
```http
GET /api/v1/places/VpOFLDRe2bzqO6Ha HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I


```
###### Response:
```http
HTTP/1.1 200 OK
Content-Length: 133
Content-Type: application/json; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT

{
	"id": "VpOFLDRe2bzqO6Ha",
	"name": "test_place_2",
	"circle": {
		"center": [
			89.000000,
			89.000000
		],
		"radius": 250
	}
}
```
----------------------------------------

## Удаление места "test_place"

###### Request:
```http
DELETE /api/v1/places/VpOFLDRe2bzqO6Ha HTTP/1.1
Host: 127.0.0.1:62441
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTI1MTAyNjAsImdyb3VwIjoidGVzdF9ncm91cCIsImlhdCI6MTQ1MjUwODQ2MCwiaWQiOiJ0ZXN0IiwiaXNzIjoiY29tLnh5enJkLmdlb3RyYWNlIiwibmFtZSI6IlRlc3QgVXNlciIsInR5cGUiOiJ1c2VyIn0.Vq2boaodeR-3jfr-y-CC_gAkzOn5FGRYqZlp-NRT07I


```
###### Response:
```http
HTTP/1.1 204 No Content
Content-Type: text/plain; charset=utf-8
Date: Mon, 11 Jan 2016 10:34:20 GMT
Content-Length: 0


```
----------------------------------------

PASS
ok  	github.com/geotrace/api	0.577s
