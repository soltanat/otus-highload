Отчет по ДЗ репликация

1. Асинхронная репликация, чтение с мастера.
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

2. Отделяем чтение в приложении

Приложение ходит в бд через pgbouncer

Результаты тестирования

Ресурсы
![Screenshot 2024-03-28 at 22.33.34.png](Screenshot%202024-03-28%20at%2022.33.34.png)

Запросы
![Screenshot 2024-03-28 at 22.33.38.png](Screenshot%202024-03-28%20at%2022.33.38.png)

