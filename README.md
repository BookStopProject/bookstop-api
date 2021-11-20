# bookstop-api

## .env

```env
DATABASE_URL=postgresql://postgres@postgres/postgres
REDIS_URL=redis://redis:6379/0
GOOGLE_API_KEY=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
API_URL=https://api.bookstop.app
APP_URL=https://bookstop.app
HMAC_SECRET=
ADMIN_AUTH=
```

## Workflow

### DB Migrate

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Generate GraphQL code

```bash
go generate ./...
```

### Generate dataloaden

```bash
go run github.com/vektah/dataloaden UserLoader string *github.com/dataloaden/example.User
```

### Build Docker

```bash
sudo docker build . -t bookstop
sudo docker tag bookstop registry.digitalocean.com/bookstop/bookstop
sudo docker push registry.digitalocean.com/bookstop/bookstop
```

### Start server

```bash
go run server.go
```

In dev, load `.env` using:

```bash
export $(grep -v '^#' .env | xargs)
```

### Start docker compose

```bash
docker compose pull
docker compose down
docker compose up
```
