# creds

This folder contains GCP service account tokens. For development and production, these tokens are optional:

* `gcp-service-account-dev.json`: For development
* `gcp-service-account-staging.json`: For end-to-end tests
* `gcp-service-account-prod.json`: For production

See `dev-scripts/generate-gcloud-auth-token` for details.
