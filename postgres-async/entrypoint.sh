#!/bin/bash

sleep 10

pg_basebackup -h postgres -D /var/lib/postgresql/data -U replication_user -P -v -X stream

touch /var/lib/postgresql/data/standby.signal

chmod 0700 /var/lib/postgresql/data

exec postgres -c hot_standby=on \
              -c primary_conninfo='host=postgres port=5432 user=replication_user password=replication_password' \
              -c shared_buffers=128kB