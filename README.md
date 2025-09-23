# INIT MVP PRODUCT

## Установка зависимостей
```bash
    go install github.com/swaggo/swag/cmd/swag@latest
```

## Docs

***- Для билда документации -***
```bash
make swagger
```


***- Для перезапуска контейнера(можно и для первого запуска) -***
```bash
make restart
```

## configs/config.json
```json
{
  "server": {
    "server_port": "1941",
    "server_mode": "development",
    "server_domain": "http://localhost:5173"
  }
}
```


## .env
```.env
MINIO_ROOT_USER=answer-minio-user
MINIO_ROOT_PASSWORD=answer-minio-password
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=answer-minio-bucket
MINIO_ENDPOINT=localhost:9000

POSTGRES_USER=answer_postgres_user
POSTGRES_PASSWORD=answer_postgres_password
POSTGRES_DB=answer_postgres_database
POSTGRES_PORT=5433
POSTGRES_HOST=localhost
POSTGRES_SSL_MODE=disable

AI_MODEL=...
AI_API_KEY=...
```
