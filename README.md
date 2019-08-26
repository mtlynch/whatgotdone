# whatgotdone

[![CircleCI](https://circleci.com/gh/mtlynch/whatgotdone.svg?style=svg&circle-token=180495ad17cc0343547e430e81d28b66ff87e9f4)](https://circleci.com/gh/mtlynch/whatgotdone)

## Architecture

### Frontend vs. Backend

TODO(mtlynch): Fill this in.

### User authentication

What Got Done uses [UserKit](https://docs.userkit.io/) for user authentication. For user signup, user login, and password reset, What Got Done loads the UserKit UI widgets in JavaScript. On the backend, the `auth` package is responsible for translating UserKit auth tokens into What Got Done usernames.

### Datastore

TODO(mtlynch): Fill this in.

### E2E tests

TODO(mtlynch): Fill this in.

## Development notes

These instructions are currently for Windows, but they'll soon be translated to bash.

### Build frontend

```
cd frontend
npm install
npm run build -- --mode development
```

### Run backend

```
cmd /c "go build --tags dev -o main.exe backend\main.go && main.exe"
```

### Run E2E tests

```
cd e2e
docker-compose up --exit-code-from cypress --abort-on-container-exit --build
```

### Quirks of the dev environment

TODO(mtlynch): Fill this in.