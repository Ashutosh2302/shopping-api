# shopping-api

First, install go and then run the following command to install the dependencies:
```bash
go mod tidy
```

Install migration tool:
```bash
brew install golang-migrate
```

Export environment varibales
```bash
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" or any postgres url
export SECRET_KEY=<any random string>
```

Run migrations
```bash
migrate -path database/migrations/ -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up
```

Finally, run the app
```bash
go run .
```