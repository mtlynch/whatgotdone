# whatgotdone

[![CircleCI](https://circleci.com/gh/mtlynch/whatgotdone.svg?style=svg&circle-token=180495ad17cc0343547e430e81d28b66ff87e9f4)](https://circleci.com/gh/mtlynch/whatgotdone) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Architecture

### Overview

![What Got Done Architecture](https://docs.google.com/drawings/d/e/2PACX-1vTolxqMjEtz6ujaM1a3ThkG3Tb1sJbv2O66TGRKVhaqNBoXtFdZjQaf3gS7l-pXbFlg02lPfM9c4foI/pub?w=917&amp;h=696)

What Got Done has a simple architecture consisting of the following parts:

* **What Got Done Backend**: A Go HTTP service (running on AppEngine) that handles all HTTP requests, datastore requests, and user authentication (via UserKit).
* **What Got Done Frontend**: A Vue2 app that renders pages in the user's browser.
* [**Cloud Firestore**](https://cloud.google.com/firestore/): What Got Done's storage provider.
* [**UserKit**](https://userkit.io): A third-party service that manages What Got Done's user authentication.

### Page Rendering Flow

What Got Done uses a somewhat unusual system for rendering pages. The Go backend first pre-renders the page server-side to populate tags related to SEO or social media that need to be set server-side. The Vue2 frontend renders the remainder of the page client-side. To avoid conflicts between the two systems' template syntax, Go uses `[[`, `]]` delimiters, while Vue uses `{{`, `}}` delimiters.

![What Got Done Render Flow](https://docs.google.com/drawings/d/e/2PACX-1vRqxoblMAAhrmI2xY_BEFmN3TRry7QdKvBOAK-1muJ79EJlJWwk1jS5t13vpjB7Kwbaf711ROMxG_cY/pub?w=1127&amp;h=1262)

The Go backend handles all of What Got Done's `/api/*` routes. These routes are What Got Done's RESTful interface between the frontend and the backend. These routes never send HTML, and instead only send JSON back and forth.

### User authentication

What Got Done uses [UserKit](https://docs.userkit.io/) for user authentication. For user signup, user login, and password reset, What Got Done loads the UserKit UI widgets in JavaScript. On the backend, the `auth` package is responsible for translating UserKit auth tokens into What Got Done usernames.

### Datastore

What Got Done uses [Google Cloud Firestore](https://firebase.google.com/docs/firestore) for data storage.

Only the What Got Done backend can access the Firestore database. Specifically, the `datastore` package manages all interactions with Firestore.

### E2E tests

What Got Done's end-to-end tests use Cypress and follow the testing pattern defined in the article [End-to-End Testing Web Apps: The Painless Way](https://mtlynch.io/painless-web-app-testing/). The testing architecture consists of two Docker containers (see [docker-compose.yml](https://github.com/mtlynch/whatgotdone/blob/master/e2e/docker-compose.yml)):

* What Got Done container
* Cypress container

The Cypress container runs a browser to exercise What Got Done's critical functionality. It uses an independent environment and credentials from the production app so that nothing in the E2E tests affect state on production UserKit or Google Cloud Platform.

To run the E2E tests yourself, see the [section below](#optional-run-e2e-tests).

## Contributing to What Got Done

Interested in contributing code or bug reports to What Got Done? That's great! Check our [Contibutor Guidelines](https://github.com/mtlynch/whatgotdone/blob/master/CONTRIBUTING.md) for more details.

## Development notes

### 0. Pre-requisites

* [Node.js](https://nodejs.org/) (8.x or higher)
* [Go](https://golang.org/dl/) (1.11 or higher)
* [Docker](https://www.docker.com/) (for E2E tests)
* [Google Cloud SDK](https://cloud.google.com/sdk/install) (for deployment)

### 1. Create a UserKit app

To run a development version of What Got Done, you'll need to create a free account on [UserKit](https://userkit.io/). Create a new UserKit app and have your UserKit App Secret Key available.

### 2. Create a project on Google Cloud Platform

What Got Done relies on Google Cloud Firestore. To run a development version of What Got Done:

1. Create a new Google Cloud Platform project.
1. Create a service account with "Owner" role.
1. Create a private key for the service account in JSON format and download it to a file called `./service-account-creds-dev.json`.

### 3. Set environment variables

Set the following environment variables:

```bash
export GOOGLE_CLOUD_PROJECT="[enter your GCP project ID]"
export USERKIT_SECRET="[enter your UserKit secret key]"
export CSRF_SECRET_SEED="any-random-string"
```

Create a file called `frontend\.env.development.local` with the following contents:

```text
VUE_APP_USERKIT_APP_ID='[your UserKit App ID]'
VUE_APP_GOOGLE_ANALYTICS_ID='0'
```

### 4. Build the frontend

To build the Vue frontend for What Got Done, run the following command:

```bash
cd frontend && \
  npm install && \
  npm run build -- --mode development
```

### 5. Run the backend

To run the Go backend server, run the following command:

```bash
mkdir bin && \
  go build --tags dev -o ./bin/main backend/main.go && \
  ./bin/main
```

What Got Done is now running on [http://localhost:3001](http://localhost:3001).

### Optional: Run frontend with hot reloading

If you're making changes to the Vue code, you'll probably want to run the standard Vue HTTP server with hot reloading. Keep the backend running, and in a separate shell session, run the following command:

```bash
cd frontend && \
  npm run serve
```

A hot-reloading Vue server will run on port [http://localhost:8085](http://localhost:8085). It will communicate with the What Got Done backend at port 3001.

#### Quirks of the dev environment

Because the production What Got Done server runs both the frontend and the backend on a single port, there are a few hacks to make a development version work:

* CORS is enabled in dev mode so that the frontend can make CORS requests to the backend from a different HTTP port.
* CSRF protection is disabled in dev mode because the Vue dev server doesn't know how to render the `<meta name="csrf-token" />` tag.
* Page titles don't work properly in dev mode because the Vue dev server doesn't know how to render the `<title>` tag.
* The Content Security Policy header in dev mode needs `unsafe-eval`, whereas we disable this in production.

### Optional: Run backend unit tests

Unit tests run in normal Golang fashion:

```bash
go test ./...
```

### Optional: Run E2E tests

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