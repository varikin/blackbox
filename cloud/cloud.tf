variable "project" {
  description = "(Required) The name of the Google Cloud Project such as `my-project-1234`"
}

variable "sheet_id" {
  description = "(Required) Google Sheet ID to append the data to"
}

variable "region" {
  default     = "us-central1"
  description = "(Optional) The region to use in Google Cloud"
}

variable "account_file" {
  default     = "account.json"
  description = "(Optional) The name of the JSON credentail file for Google Cloud"
}

variable "function_name" {
  default     = "black-box-stream-to-sheet"
  description = "(Optional) The name of the Google Cloud Function to create"
}

variable "topic_name" {
  default     = "black-box-stream"
  description = "(Optional) The name of the Google Pub/Sub Topic to create"
}

variable "bucket_name" {
  default     = "black-box-functions"
  description = "(Optional) The name of the Google Storage bucket to create"
}

# Use Google Cloud provider with an account file in the current directory
provider "google" {
  credentials = "${file(var.account_file)}"
  project     = "${var.project}"
  region      = "${var.region}"
}

# Load the Archive provider
provider "archive" {}

# Create a Pub/Sub Topic
resource "google_pubsub_topic" "topic" {
  name = "${var.topic_name}"
}

# Add the Particle service account as a Pub/Sub publisher to the topic
resource "google_pubsub_topic_iam_binding" "publisher" {
  topic = "${google_pubsub_topic.topic.name}"
  role  = "roles/pubsub.publisher"

  members = [
    "serviceAccount:particle-public@particle-public.iam.gserviceaccount.com",
  ]
}

# Create a bucket to store the Function in
resource "google_storage_bucket" "bucket" {
  name = "${var.bucket_name}"
}

# Zip up the local Function so it can be put in the bucket
data "archive_file" "function_bundle" {
  type        = "zip"
  source_dir  = "./src"
  output_path = "dist/${var.function_name}.zip"
}

# Put the Function zip in the bucket
# Adding the sha256 in the name to cache-bust the Function loader
resource "google_storage_bucket_object" "function_object" {
  name   = "${var.function_name}-${data.archive_file.function_bundle.output_base64sha256}.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "${data.archive_file.function_bundle.output_path}"
}

# Create the Function with the Go 1.11 runtime
# The Function is triggered by the topic
resource "google_cloudfunctions_function" "function" {
  name                  = "${var.function_name}"
  available_memory_mb   = 128
  timeout               = 60
  runtime               = "go111"
  source_archive_bucket = "${google_storage_bucket_object.function_object.bucket}"
  source_archive_object = "${google_storage_bucket_object.function_object.name}"

  event_trigger = {
    event_type = "google.pubsub.topic.publish"
    resource   = "${google_pubsub_topic.topic.name}"
  }
  environment_variables = {
    SHEET_ID = "${var.sheet_id}"
  }

  entry_point = "Run"
}
