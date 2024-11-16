# Scaffold Echo Application in GO

## Project Dependencies
- [Go](https://go.dev/doc/install)
- [NodeJS](https://nodejs.org/en/download/prebuilt-installer/current) (For tailwind)
- Echo
- [Templ](https://templ.guide/quick-start/installation)
- [Air](https://github.com/air-verse/air) (Only Dev)
- [TailwindCSS](https://tailwindcss.com/)
- Daisy UI
- [SQLC](https://docs.sqlc.dev/en/latest/overview/install.html)
- [GoMigrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- PostgreSQL

All these dependencies are basically building tools to make the project work

## To start the project
First check if all dependencies are installed by running. You can check the Makefile to alter the values with the names you want.

```shell
go get .
npm i
make build
make create-db
air
```

## Development Tips
### How to handle migrations
We're using [GoMigrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate), read the documentation for more information.

> The lib allows to use an ENV Var to streamline the process so you don't need to always put the database URL in the commands
> just set POSTGRESQL_URL in your shell env and you're good to go

** Here's a fast command to create a new migration for the project: **

```shell
# This URL is pointing to the DB created with the variables set in Makefile, please alter it to your database name and credentials
# It's also good to only use this in development, since the string will contain the db password
migrate create -ext sql -dir ./migrations -seq name_of_migration -database 'postgres://postgres:password@localhost:5432/echo_scaffold?sslmode=disable'
```

** Run the migrations **
```shell
migrate -database 'postgres://postgres:password@localhost:5432/echo_scaffold?sslmode=disable' -path ./migrations up
```


## To build for production
