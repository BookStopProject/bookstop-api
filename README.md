# bookstop-api

## Workflow

### DB Migrate

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Generate GraphQL code

```bash
go generate ./...
```

### Start server

```bash
go run server.go
```

In dev, load `.env` using:

```bash
export $(grep -v '^#' .env | xargs)
```
