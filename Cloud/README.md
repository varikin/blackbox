# Google Cloud Integration

Requires a Google Cloud project.

## Terraform

Requires version 0.11

Note that Terraform 0.12 is not supported due to configuration language changes.

## Service Account

A service account is required with a JSON based key.
Permissions required:

* Pub/Sub Admin
* Cloud Functions Developer
* Storage Admin
* Service Account User

Download the key as a JSON file and place in this directory named `account.json`.