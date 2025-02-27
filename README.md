# News API

```
git pull git@github.com:Dima-Karpov/comment-api.git
make install
```

### migrate
    1. docker-compose exec db bash
    2. psql -h db -U postgres -d postgres -f /migrations/schema/000001_init.up.sql (накатываем миграции)
    3. psql -U postgres -d api
    4. \dt (увидеть все таблицы)

### не забыть создать общую сеть если еще не создана
    docker network create my_custom_network