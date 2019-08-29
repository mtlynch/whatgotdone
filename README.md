# whatgotdone

[![CircleCI](https://circleci.com/gh/mtlynch/whatgotdone.svg?style=svg&circle-token=180495ad17cc0343547e430e81d28b66ff87e9f4)](https://circleci.com/gh/mtlynch/whatgotdone)

## Architecture

### Frontend vs. Backend

TODO(mtlynch): Fill this in.

### User authentication

What Got Done uses [UserKit](https://docs.userkit.io/) for user authentication. For user signup, user login, and password reset, What Got Done loads the UserKit UI widgets in JavaScript. On the backend, the `auth` package is responsible for translating UserKit auth tokens into What Got Done usernames.

### Datastore

What Got Done uses [Google Cloud Firestore](https://firebase.google.com/docs/firestore) for data storage.

Only the What Got Done backend can access the Firestore database. Specifically, the `datastore` package manages all interactions with Firestore.

### E2E tests

TODO(mtlynch): Fill this in.

## Development notes

### Pre-requisites

* [Node.js](https://nodejs.org/) (8.x or higher)
* [Go](https://golang.org/dl/) (1.11 or higher)
* [Docker](https://www.docker.com/) (for E2E tests)
* [Google Cloud SDK](https://cloud.google.com/sdk/install) (for deployment)

In addition, you must create a [UserKit](https://userkit.io/) account and have your UserKit App Secret Key available.

### Set environment variables

```bash
export GOOGLE_CLOUD_PROJECT="[enter your GCP project ID]"
export USERKIT_SECRET="[enter your UserKit secret key]"
```

### Build frontend

```bash
cd frontend && \
  npm install && \
  npm run build -- --mode development
```

### Run backend

```bash
mkdir bin && \
  go build --tags dev -o ./bin/main backend/main.go && \
  ./bin/main
```

### Run E2E tests

To run the end to end tests, you'll need to create a dedicated GCP project. You can reuse your dev project, but the E2E tests will write to the datastore for the GCP project you specify. Specify the GCP project in `e2e\docker-compose.yml` under `GOOGLE_CLOUD_PROJECT`.

You'll also need to create a dedicated UserKit app. Save the UserKit app secret in a file called `e2e/staging-secrets.env` with the following contents:

```text
USERKIT_SECRET=[your userkit app secret]
```

You'll also need to manually create a user for that app with the following credentials:

* Username: `staging.jimmy`
* Password: `just4st@ginG!`

When you've completed these steps, you can run the E2E tests as follows:

```bash
cd e2e && \
  docker-compose up --exit-code-from cypress --abort-on-container-exit --build
```

### Quirks of the dev environment

TODO(mtlynch): Fill this in.