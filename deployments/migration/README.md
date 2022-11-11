# Database migration

Use [sql-migrate](https://github.com/rubenv/sql-migrate) to migrate SQL schema defined in .sql file.

# Installation

To install sql-migrate, use the following:

```bash
go get -v github.com/rubenv/sql-migrate/...
```

For Go version from 1.18, use:

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

# Usage

Modify mysql config in dbconfig.yml, or setup environment variables, such as DB_USER used by dbconfig.yml, then exec following commands:

```bash
# migrate up
make migrate-up

# migrate status
make migrate-status
```

<details>
    <summary>Migration commands available</summary>

| Command             | Desc                                           |
| ------------------- | ---------------------------------------------- |
| `make migrate-up`   | runs migration up command                      |
| `make migrate-down` | runs migration down command                    |
| `make force`        | Set particular version but don't run migration |
| `make goto`         | Migrate to particular version                  |
| `make drop`         | Drop everything inside database                |
| `make create`       | Create new migration file(up & down)           |

</details>
