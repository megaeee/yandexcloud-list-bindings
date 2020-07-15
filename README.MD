# Yandex.Cloud. Serverless-функция для просмотра всех ресурсов, доступных пользователю
Данная функция обходит все сущности, для которых доступно назначение ролей в облаке, и проверяет наличие пользователя в биндингах

## Подготовка окружения
1. [Создать сервисный аккаунт](https://cloud.yandex.ru/docs/iam/operations/sa/create)
1. [Выдать права доступа к сервисному аккаунту](https://cloud.yandex.ru/docs/iam/operations/sa/set-access-bindings)
1. [Выдать права доступа к облаку для сервисного аккаунта](https://cloud.yandex.ru/docs/resource-manager/operations/cloud/set-access-bindings#access-to-sa) (Роль *viewer* на облако)
1. [Создать функцию](https://cloud.yandex.ru/docs/functions/quickstart/function-quickstart). Для запуска кода из репозитория необходимо сохранить его содержимое в zip-архив и передать в форме редактирования после указания соответствующего метода загрузки, а также указать следующие параметры: 
  * *Среда выполнения* - golang114
  * *Таймаут, c* - 300. Такой большой интервал нужен потому, что функция обходит все возможные сущности в пределах пользовательского облака
  * *Сервисный аккаунт* - созданный ранее сервисный аккаунт
## Запуск и тестирование:
    curl -s -XPOST https://functions.yandexcloud.net/{function_id}?integration=raw -d '{"cloudId": "{cloud_id}", "user": "{user_id_or_name}"' -H "Authorization: Bearer $(yc iam create-token)"
json-параметры для запуска:
1. *cloud_id* - идентификатор облака, в пределах которого ищем доступы
1. *user_id_or_name* - идентификатор или имя пользователя, чьи доступы ищем

Выводом функции будет вывод формата, указанного ниже, либо Golang-ошибка

    {
        "statusCode": 200,
        "body": [
            {
            "binding": {
                "role_id": "resource-manager.clouds.member",
                "subject": {
                "id": "b1g5q8h52kmbindingid",
                "type": "userAccount"
                }
            },
            "resourceType": "cloud",
            "resourceID": "aje46hba9nluocloudid"
            },
            ...
        ]
    }
