# URL_Shortener Service

## Описание

Rest Api сервис для укорачивания ссылок.

- **Технологии:**
    - Go
    - Docker
    - Chi Router
    - JWT
    - PostgreSQL

- **Эндпоинты:**
    - `POST /` - сохранить ссылку
    - `DELETE /{alias}` - удалить ссылку
    - `GET /{alias}}` - редирект по укороченной ссылке
