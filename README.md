# bookstop-api

API Server for BookStop, written in Golang. Exchange books with each other, share your thoughts, participate in events, all in one stop.

Website: https://bookstop.app

| Codebase                                                            |                Description                 |
| :------------------------------------------------------------------ | :----------------------------------------: |
| [bookstop-front](https://github.com/BookStopProject/bookstop-front) |        Next.js Frontend Application        |
| [bookstop-api](https://github.com/BookStopProject/bookstop-api)     | Golang API Server (GraphQL) and Admin Page |
| [bookstop-home](https://github.com/BookStopProject/bookstop-home)   |              Preact homepage               |

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
ADMIN_AUTH=username:password
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

### Start server (dev)

Start docker compose:

```bash
docker compose -f docker-compose.dev.yml up
```

Run server:

```bash
go run server.go
```

### Start docker compose

```bash
docker compose down
docker compose up --build -d
```
