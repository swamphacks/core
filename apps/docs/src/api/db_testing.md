# Database testing

## Manual

1. Make sure the docker instance is currently running (i.e. you ran `$ docker compose up`)
!!! note

    You may have to run docker with sudo depending on your system configuration.

1. In a seperate terminal (I suggest you try using tmux) list the current docker processes
``` bash
docker ps
```

1. Copy the process id for the docker container running postgres, and paste into the following command in order to spawn a shell with access to the container.
``` bash
sudo docker exec -it ef7XXXXXX07e sh
```

1. Connect to the postgres database using a database url. The url can be found inside `core/apps/api/.env.example`. 
``` bash
psql postgres://postgres:postgres@postgres:5432/coredb
```

1. Check to see if database tables currently exist. You can by listing the currently created tables.
``` bash
\dt
```
If there is nothing here, then go into `/core/apps/api` and run `make migrate`, which runs sql commands added to [migrations](https://en.wikipedia.org/wiki/Schema_migration) in `core/apps/api/internal/db/migrations`.
1. You can now test to see if rows, columns and tables are updated appropriately with psql commands. Use a [reference to psql](https://www.postgresql.org/docs/17/app-psql.html) if you need help finding commands.
