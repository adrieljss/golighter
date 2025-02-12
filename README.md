![GoLighter Logo](https://github.com/user-attachments/assets/b23ad30b-b450-48d1-a9fb-2a802e5f6e72)

# GoLighter
GoLighter is a lightweight API boilerplate, built with a user system and authentication.

## Features

- User management
- Authentication
- Authorization
- Flag-based access control
- JWT-based authentication
- PostgreSQL Database
- Redis Caching

## Getting Started
```bash
go run main.go
```
and then run the migrations in `migrations/` folder.

## Testing

To run tests, make sure to activate the testing database first (in a seperate db from production, and already have migrations). Run the following command to recursively run tests.

```bash
go run ./...
# or with a verbose output
go run -v ./...
```

Or additionally, you can use [tparse](https://github.com/mfridman/tparse/):
```bash
set -o pipefail && go test -json ./... | tparse -all # linux
go test -json ./... | tparse -all # windows
```