
## Сервис авторизации

### API
 1. Регистрация 

    - Метод - `POST`
    - Путь - `/register` 
    - Если добавить query-параметр `role=manager` или `role=chef`, то новому пользователю будет выдана эта роль 
    - Сообщение - ```{"username": "<username>", "email" : "<email>", "password": "<password>"}```
    - Пример curl-запроса:
      ```sh
      curl localhost:8080/register?role=manager --data '{"username": "ivan", "email" : "ivan@mail.ru", "password": "ivan"}' -vv
      ```
    - Пример ответа - `{"message":"success"}`
    - Возможные статус коды:
        - 200 - все ок, возвращается пользователь
        - 400 - не валидный ответ
        - 409 - уже существует пользовательс с таким username/email 
        - 500 - ошибки на бекенде
    
  2. Авторизация
     - Метод - `POST`
     - Путь - `/login`
     - Сообщение - ```{"email" : "<email>", "password": "<password>"}```
     - Пример curl-запроса:
       ```sh
       curl localhost:8080/login --data '{"email" : "ivan1@mail.ru", "password": "ivan"}' -vv
       ```
     - Пример ответа - `{"jwt_token":"20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2"}`
     - Возможные статус коды:
       - 200 - все ок, возвращается токен
       - 400 - ошибки не валидная почта/пароль
  3. Получение пользователя
     - Метод - `GET`
     - Путь - `/user`
     - Требуется заголовок `"Authorization: Bearer <jwt_token>"`, где jwt_token получается из запроса авторизации
         - Пример curl-запроса:
       ```sh
       curl localhost:8080/user -H 'Authorization: Bearer 20ee6dac6bd03313ee389ac566e0426afc89a752de8655ca14c04f39d76eb7e2' -vv
       ```
     - Пример ответа 
       ```
       {
           "id": 20,
           "username": "ivan",
           "email": "ivan@mail.ru",
           "role": "customer"
       }
       ```
     - Возможные статус коды:
       - 401 - предоставлен протухший/несуществующий токен
       - 200 - все ок, возвращается пользователь
       - 500 - ошибки на бекенде
