# go project layout

This is a simple clean architecture template with gin framework, fx as dependency container, gorm as orm for database related operations.

## Implemented Features

- Dependency Injection (go-fx)
- Routing (gin web framework)
- Logging provides context trace(trace_id) and log file rotation ([zap](https://github.com/uber-go/zap) and [lumberjack](https://github.com/natefinch/lumberjack))
- Gin middlewares setup (cors, zap logger with custom log format)
- Database Setup (mysql)
- Models Setup and Automigrate (gorm with zap logger)
- Authentication (JWT)
- Migration Runner Implementation
- Cobra Commander CLI Support. try: `go run . --help`
- Pprof setup with gin

## Running the project

- Make sure you have [go](https://go.dev/) and git installed.
- Clone and build app:
  ```bash
  $ go clone https://github.com/dean2032/go-project-layout.git
  $ cd go-project-layout
  $ go build -o app
  ```
- Run a subcommand of app. Two subcommands are provided:
  - file_server: this subcommand will run a file download web server exposing the current directory on port 8888:
    ```bash
    $ ./app file_server -verbose
    ```
  - api_server: this subcommand will run a http api server. Mysql db must be installed and configured correctly before it can run properly. See [./deployments/migration](https://github.com/dean2032/go-project-layout/tree/main/deployments/migration) for more details.

Have fun!
