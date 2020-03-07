# backup-service

## Overview

This is a daily backup service for What Got Done's Filestore data. It exports all Filestore entities to the What Got Done backup GCS bucket and is set up to run every 24 hours according to an AppEngine cron job.

## Source

Most of the code comes from the GCP tutorial: https://firebase.google.com/docs/firestore/solutions/schedule-export

## Pre-requisites

For this to work, AppEngine service worker must have correct privileges

```bash
PROJECT_ID=whatgotdone

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member "serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com" \
    --role roles/datastore.importExportAdmin
```

## Deploying

This is not set up to auto-deploy with the regular site because deployments are required infrequently.

To deploy this service manually:

```bash
gcloud app deploy app.yaml
```
