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

Результат тестирования.
Ресурсы
![Screenshot 2024-03-28 at 20.31.28.png](Screenshot%202024-03-28%20at%2020.31.28.png)
Запросы
![Screenshot 2024-03-28 at 20.31.31.png](Screenshot%202024-03-28%20at%2020.31.31.png)


Мониторинга контейнеров с помощью cadvisor, не нашел как получить статистику la, она в нуле, возможно из-за особенностей mac m1.

2. Отделяем чтение в приложении

