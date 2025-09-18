# INIT MVP PRODUCT

#### Установка зависимостей
```bash
    go install github.com/swaggo/swag/cmd/swag@latest
```

#### Docs
```bash
make swagger
```

### configs/config.json
```json
{
  "server": {
    "server_port": "1941",
    "server_mode": "development",
    "server_domain": "http://localhost:5173"
  },
  "business-database": {
    "db_password":"answer_postgres_password",
    "db_host": "localhost",
    "db_port": "5433",
    "db_username":"answer_postgres_user",
    "db_name":"answer_postgres_database",
    "db_ssl_mode": "disable"
  },
  "ai": {
    "model": "gemini-you",
    "api_key": "keeeey"
  }
}
```
