Отчет по ДЗ репликация

## 1. Асинхронная репликация, чтение с мастера.
Настройки мастера, для наглядности уменьшен `shared_buffer=128kB`.

```shell
postgres
    -c listen_addresses='*'
    -c wal_level=hot_standby
    -c hba_file=/etc/postgresql/pg_hba.conf
    -c log_statement=all
    -c shared_buffers=128kB
```

Настройка слейвов
```shell
postgres 
    -c hot_standby=on
    -c primary_conninfo='host=postgres port=5432 user=replication_user password=replication_password'
    -c shared_buffers=128kB
```

Результат тестирования.

Ресурсы
![Screenshot 2024-03-28 at 20.31.28.png](Screenshot%202024-03-28%20at%2020.31.28.png)

Запросы
![Screenshot 2024-03-28 at 20.31.31.png](Screenshot%202024-03-28%20at%2020.31.31.png)


Мониторинг контейнеров с помощью cadvisor, не нашел как получить статистику la, она в нуле, возможно из-за особенностей mac m1.


## 2. Отделяем чтение в приложении

Приложение ходит в бд через pgbouncer

Результаты тестирования

Ресурсы
![Screenshot 2024-03-28 at 22.33.34.png](Screenshot%202024-03-28%20at%2022.33.34.png)

Запросы
![Screenshot 2024-03-28 at 22.33.38.png](Screenshot%202024-03-28%20at%2022.33.38.png)

По графикам видно, что нагрузка теперь на слейвы


## 3. Синхронная репликация, потеря транзакций

Запускаем [postgres-async-docker-compose.yaml](..%2Fpostgres-async-docker-compose.yaml)

Для того что бы запустить синхронную репликацию, нужно добавить к команде запуска или в конфиг параметр `-c synchronous_standby_names='replica0,replica1'`

После того как включили синхронную репликацию, запускаем скрипт генерирующий вставки в тестовую таблицу

```shell
make run-gentransactions
```

Во время выполнения скрипта убиваем мастер docker kill postgres-master и останавливаем скрипт он выдаст нам кол-во выполненных вставок: 608

Сравним кол-во с кол-вом строк в slave'аx 
```sql
SELECT COUNT(*) FROM public.test_transactions t
```
В обоих слейвах так же 608 строк.

Подключимся к слейву: `docker exec -it postgres-slave-0 bash`, промоутим до мастера `pg_ctl promote -D /var/lib/postgresql/data/`

Перезапускаем postgres-replica-1 с новым primary_conninfo, находится в postgres-sync.

Теперь слейв 1 синхронизируется со слейв 0