
## Сервис заказов

#### Стоит знать:
Запускает обработчик заказов, который работает в фоне от апи. Он каждые 10 секунд ищет заказ в статусе waiting, если находит, то переводит в статус processing. Потом проверяет наличие блюд требуемых для заказа -- если не хватает, то переводит заказ в статус cancel. Если блюд хватает, то он выполняет заказ в течение 30 секунд и в конце переводит в статус success. Дальше цикл начинается с начала. 

### API
1. Получение меню 
    - _Менеджер видит все блюда, Покупатель только доступные(is_available=true)_
    - Метод - `GET`
    - Путь - `/dishes`
    - Требуется заголовок (опционально) `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
        - Пример curl-запроса:
      ```sh
      curl localhost:8081/dishes -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' -vv
      ```
    - Пример ответа 
      ```
        {
          "dishes": [
            {
              "id": 1,
              "name": "Pizza",
              "description": "with pineapples",
              "price": 10,
              "quantity": 20,
              "is_available": true
            },
            {
              "id": 2,
              "name": "Cakes",
              "description": "with pineapples",
              "price": 5,
              "quantity": 10,
              "is_available": true
            }
          ]
        }
      ```
    - Возможные статус коды:
        - 200 - все ок, возвращается пользователь
        - 500 - ошибки на бекенде

2. Создание заказов
   - Метод - `POST`
   - Путь - `/orders`
   - Требуется заголовок `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
   - Тело запроса ```{"dishes": {dish_id: queantity, dish_id2: queantity2}, "special_requests": "<special_requests>"}``` 
   - Пример curl-запроса:
     ```sh
     curl localhost:8081/orders -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' \
     --data '{"dishes": {"1": 1, "2": 1}, "special_requests": "поострее, пожалуйста"}' -vv
     ```
   - Пример ответа 
      ```
     {"order_id": 1}
      ```
   - Возможные статус коды:
      - 200 - все ок, возвращается id заказа
      - 400 - заказаны несуществующие блюда/нехватает ресурсов
      - 401 - пользователь не авторизован
      - 500 - ошибки на бекенде

3. Получение статуса заказа
   - Метод - `GET`
   - Путь - `/orders/{order_id}`
   - Пример curl-запроса:
     ```sh
     curl localhost:8081/orders/1 -vv
     ```
   - Пример ответа
      ```
           {
        "id": 4,
        "user_id": 20,
        "status": "success",
        "dishes": [
          {
            "dish_id": 1,
            "quantity": 1,
            "price": 10
          },
          {
            "dish_id": 2,
            "quantity": 1,
            "price": 10
          }
        ],
        "special_requests": "поострее, пожалуйста"
      }
      ```
   - Возможные статус коды:
      - 200 - все ок, возвращается заказ и его статус
      - 400 - неправильный путь
      - 500 - ошибки на бекенде

4. Добавление блюд
    - Метод - `POST`
    - Путь - `/dishes`
    - Требуется заголовок `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
    - Тело запроса ```{"name": "<name>", "description": "<description>", "price": price, "quantity": quantity, "is_available": true/false}```
    - Пример curl-запроса:
      ```sh
      curl localhost:8081/dishes  -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' \
      --data '{"name": "Pizza", "description": "with pineapples", "price": 10, "quantity": 20, "is_available": true}' -vv
      ```
    - Возможные статус коды:
        - 200 - все ок, блюдо добавлено
        - 405 - нет прав

5. Правки информации о блюдах
    - Метод - `PUT`
    - Путь - `/dishes/{dishID}`
    - Требуется заголовок `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
    - Тело запроса ```{"name": "<name>", "description": "<description>", "price": price, "quantity": quantity, "is_available": true/false}```
    - Пример curl-запроса:
      ```sh
      curl -X PUT localhost:8081/dishes/3  -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' \
      --data '{"name": "Pizza", "description": "with pineapples", "price": 20, "quantity": 30, "is_available": true}' -vv
      ```
    - Возможные статус коды:
        - 200 - все ок, блюдо добавлено
        - 405 - нет прав

6. Удаление блюд
    - Метод - `DELETE`
    - Путь - `/dishes/{dishID}`
    - Требуется заголовок `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
    - Пример curl-запроса:
      ```sh
      curl -X DELETE localhost:8081/dishes/3  -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' -vv
      ```
    - Возможные статус коды:
        - 200 - все ок, блюдо удалено
        - 405 - нет прав
