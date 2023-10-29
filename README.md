### Запуск

```
go mod tidy
go build -o app_linux
./app_linux
```


### Структура каталога

    /conf
        app.conf Конфигурации
    /routers
        /routeHelpers   Middleware и пр. вспомогательные функции роутинга
        /routes         Маршруты
        router.go       Файл инициализации маршрутизатора
    /src
        /components     Компоненты приложения
            /auth
            /forms      Типы структур для данных из query params
            /requests   Типы структур для данных из body json
            /responses  Типы структур ответа и вспомогательные функции
        /controllers    MVC контроллеры
        /models         MVC модели
    /swagger            Статические файлы swagger
    main.go             Главнюк
        

### База данных

Не стал тащить сюда утилиту фреймворка для миграций, так как структура данных слишком простая, выгрузил схему в файл ``db.sql``

Для подклюения в файле ``conf/app.conf`` настроить:

```
[database]
driver = postgres
user = www
password = secret
host = localhost
port = 5432
name = db
ssl = disable
```

### Swagger

http://localhost:8080/swagger/