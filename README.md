# whatgotdone

[![CircleCI](https://circleci.com/gh/mtlynch/whatgotdone.svg?style=svg&circle-token=180495ad17cc0343547e430e81d28b66ff87e9f4)](https://circleci.com/gh/mtlynch/whatgotdone) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Architecture

### Overview

![What Got Done Architecture](https://docs.google.com/drawings/d/e/2PACX-1vTolxqMjEtz6ujaM1a3ThkG3Tb1sJbv2O66TGRKVhaqNBoXtFdZjQaf3gS7l-pXbFlg02lPfM9c4foI/pub?w=917&amp;h=696)

What Got Done has a simple architecture consisting of the following parts:

* **What Got Done Backend**: A Go HTTP service that handles all HTTP requests, datastore requests, and user authentication (via UserKit).
* **What Got Done Frontend**: A Vue2 app that renders pages in the user's browser.
* [**SQLite**](https://www.sqlite.org/): What Got Done's storage provider.
* [**Litestream**](https://litestream.io): (optional) Syncs What Got Done's SQLite database to cloud storage.
* [**UserKit**](https://userkit.io): A third-party service that manages What Got Done's user authentication.

### Page Rendering Flow

What Got Done uses a somewhat unusual system for rendering pages. The Go backend first pre-renders the page server-side to populate tags related to SEO or social media that need to be set server-side. The Vue2 frontend renders the remainder of the page client-side. To avoid conflicts between the two systems' template syntax, Go uses `[[`, `]]` delimiters, while Vue uses `{{`, `}}` delimiters.

![What Got Done Render Flow](https://docs.google.com/drawings/d/e/2PACX-1vRqxoblMAAhrmI2xY_BEFmN3TRry7QdKvBOAK-1muJ79EJlJWwk1jS5t13vpjB7Kwbaf711ROMxG_cY/pub?w=1127&amp;h=1262)

The Go backend handles all of What Got Done's `/api/*` routes. These routes are What Got Done's RESTful interface between the frontend and the backend. These routes never send HTML, and instead only send JSON back and forth.

### User authentication

What Got Done uses [UserKit](https://docs.userkit.io/) for user authentication. For user signup, user login, and password reset, What Got Done loads the UserKit UI widgets in JavaScript. On the backend, the `auth` package is responsible for translating UserKit auth tokens into What Got Done usernames.

### Integration tests

What Got Done's integration tests use Cypress and follow the testing pattern defined in the article [End-to-End Testing Web Apps: The Painless Way](https://mtlynch.io/painless-web-app-testing/). The testing architecture consists of four Docker containers (see [docker-compose.yml](https://github.com/mtlynch/whatgotdone/blob/master/integration/docker-compose.yml)):

* What Got Done container
* Cypress container
* Test data manager

The integration test is entirely self-contained, so nothing that happens during testing affects anything outside of the docker-compose environment.

To run the integration tests yourself, see the [section below](#optional-run-integration-tests).

## Contributing to What Got Done

Interested in contributing code or bug reports to What Got Done? That's great! Check our [Contibutor Guidelines](https://github.com/mtlynch/whatgotdone/blob/master/CONTRIBUTING.md) for more details.

## QuickStart

To run What Got Done in a Docker container, run

```bash
docker-compose up
```

What Got Done will be running at [http://localhost:3001](http://localhost:3001).

Dev-mode authentication uses [UserKit dummy mode](https://docs.userkit.io/docs/dummy-mode). You can log in with any username using the password `password`.

## Development notes

### Pre-requisites

* [Node.js](https://nodejs.org/) (12.x or higher)
* [Go](https://golang.org/dl/) (1.17 or higher)
* [Docker](https://www.docker.com/) (for E2E tests)
* [Google Cloud SDK](https://cloud.google.com/sdk)
* [screen](https://wiki.debian.org/screen)

### Start hot reloading server

Run the following command to start a What Got Done development server:

```bash
./dev-scripts/serve
```

1. Builds the Vue frontend.
1. Starts a hot reloading server for the Vue frontend.
1. Starts a hot reloading server for the backend.

* The backend server will run on [http://localhost:3001](http://localhost:3001).
* The frontend server will run on port [http://localhost:8085](http://localhost:8085).

#### Quirks of the dev environment

Because the production What Got Done server runs both the frontend and the backend on a single port, there are a few hacks to make a development version work:

* CORS is enabled in dev mode so that the frontend can make CORS requests to the backend from a different HTTP port.
* CSRF protection is disabled in dev mode because the Vue dev server doesn't know how to render the `<meta name="csrf-token" />` tag.
* Page titles don't work properly in dev mode because the Vue dev server doesn't know how to render the `<title>` tag.
* The Content Security Policy header in dev mode needs `unsafe-eval` and `unsafe-inline`, whereas we disable this in production.

### Optional: Run backend unit tests

Unit tests run in normal Golang fashion:

```bash
go test ./...
```

### Optional: Run integration tests

Integration tests run all components together using [UserKit dummy mode](https://docs.userkit.io/docs/dummy-mode) as authentication:

```bash
dev-scripts/run-integration-tests
```

### Optional: Enable public analytics from Google Analytics

What Got Done supports pulling metrics from Google Analytics into the page content. To enable this:

1. Enable the [Google Analytics Reporting API](https://console.cloud.google.com/apis/library/analyticsreporting.googleapis.com) in your Google Cloud Platform project.
1. Create a [service account](https://console.cloud.google.com/iam-admin/serviceaccounts) in Google Cloud Platform console for your What Got Done project.
   1. Assign the service account no permissions/roles, but save its private key as JSON.
   1. Click "Create Key" to create a private key and save it in JSON format as `gcp-service-account-prod.json` in the What Got Done root directory.
1. In Google Analytics, open Admin > View > View User Management and add the email address of the service account you just created (it will have an email like `[name]@[project ID].iam.gserviceaccount.com`.
   1. Grant the user only "Read & Analyze" permissions.
1. In Google Analytics, open Admin > View > View Settings
   1. Save the View ID as an environment variable like `export GOOGLE_ANALYTICS_VIEW_ID=12345789`

### Optional: Enable image uploads

What Got Done optionally allows image uploads from users. To enable this:

1. Create a Google Cloud Storage bucket
1. Choose uniform permissions for the bucket
1. Add to the bucket permissions for the `allUsers` user with the role "Storage Object Viewer"
  * This makes all images in the bucket world-readable.
1. When launching What Got Done, set the environment variable `PUBLIC_GCS_BUCKET` to the name of your GCS bucket.

When users paste images into their What Got Done entries, they will upload to your GCS bucket and auto-link from the entry editor.
