- [Setup](#setup)
  - [Docker](#docker)
- [Hello World](#hello-world)
  - [docker-compose](#docker-compose)
  - [db](#db-mysql)
  - [migration](#migration-flyway)
  - [api](#api-golang)
- [Bye World](#bye-world)

# Setup

### Docker

- Docker Desktop
  - https://docs.docker.com/desktop/
- docker-compose 3.8に必要な
  - Docker Engine: [19.03.0+](https://docs.docker.com/compose/compose-file/)
  - docker-compose : [1.25.5以上](https://docs.docker.com/compose/release-notes/#1255)

```
docker -v
Docker version 19.03.8, build afacb8)

docker-compose --version
docker-compose version 1.25.5, build 8a1c60f6
```

# Hello World

### docker-compose
```
> make docker-compose/build
> make docker-compose/up
```

### db (mysql)

データベースの新規作成。
```
> make mysql/init
> make mysql/client

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| ...                |
| table       |
+--------------------+
5 rows in set (0.01 sec)

```

### migration (flyway)

flyway: https://flywaydb.org/documentation/

#### baseline
```
// 初期化コマンド
> make flyway/baseline
```

#### migrate

```
> make flyway/info
...

> make flyway/migrate
...
> make flyway/info
+-----------+---------+-----------------------+----------+---------------------+----------+
| Category  | Version | Description           | Type     | Installed On        | State    |
+-----------+---------+-----------------------+----------+---------------------+----------+
|           | 1       | << Flyway Baseline >> | BASELINE | 2020-07-13 02:24:32 | Baseline |
| Versioned | 2       | create table user     | SQL      | 2020-07-13 02:26:00 | Success  |
| Versioned | 3       | create table room     | SQL      | 2020-07-13 02:26:01 | Success  |
+-----------+---------+-----------------------+----------+---------------------+----------+
```

### api (golang)

```bash
❯ curl -v http://localhost:10001/todo/login | jq .
< HTTP/1.1 400 Bad Request
< Content-Type: application/json
{
  "message": "authorization header not found"
}
```

authorization headerないから認証エラーになる。（正しい挙動)
```

# Bye World
```
make docker-compose down
```

### Bye db
```
make __drop/mysql
```
