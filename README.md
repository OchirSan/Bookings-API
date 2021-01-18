# Booking API


###### **Get /metrics**

Метрики


###### **GET /rooms/list**

Выводит комнату по room_id, если room_id не указан, то выводит все комнаты <br>
Пример запроса : http://localhost:9000/rooms/list?room_id=1
: <br>
```json
{
    "room_id": 1
    "description":"Poors man room",
    "price":1
}
```

###### **GET /bookings/list**

Выводит брони по указанному room_id<br>
пример запроса : http://localhost:9000/bookings/list?room_id=1<br>
```json
[
  {
    "user_id":1,
    "comment":"Начисление средств",
    "amount":23000,
    "date":"2020-09-27T15:06:03.137484+03:00"
    },
  {
    "user_id":1,
    "comment":"Начисление средств",
    "amount":11000,
    "date":"2020-09-27T00:23:08.527134+03:00"
  },
  {
    "user_id":1,
    "comment":"Перевод средств",
    "amount":11000, 
    "date":"2020-09-27T00:23:08.527134+03:00"
  }
]
```


###### **POST /api/v1/balance**

Создает пользователя с балансом  <br>
на входе:
 ```json
  {
    "user_id":1, 
    "amount":10000
  }
```


###### **PUT /api/v1/accrual**

Начисляет пользователю указанную сумму, а также записывает транзакцию  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **PUT /api/v1/debit**

Списывает у пользователя указанную сумму, а также записывает транзакцию  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **PUT /api/v1/transfer/{id:[0-9]+}**

Списывает у пользователя в теле запроса и начисляет пользователю в урле, также записывает транзакции  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **DELETE /api/v1/balance/{id:[0-9]+}**

Удаляет пользователя по id  <br>




  