variable "project" {
  description = "(Required) The name of the Google Cloud Project such as `my-project-1234`"
}

variable "sheet_id" {
  description = "(Required) Google Sheet ID to append the data to"
}

variable "device_ids" {
    description = "(Required) A comma delimited list of Particle device ids that have the sensors"
}

variable "particle_access_token" {
  description = "(Required) The access token for the Particle API"
}

variable "region" {
  default     = "us-central1"
  description = "(Optional) The region to use in Google Cloud"
}

variable "account_file" {
  default     = "account.json"
  description = "(Optional) The name of the JSON credentail file for Google Cloud"
}

variable "sensor_data_function" {
  default     = "black-box-stream-to-sheet"
  description = "(Optional) The name of the Google Cloud Function to handle sensor data"
}

variable "sensor_data_topic" {
  default     = "black-box-sensor-data"
  description = "(Optional) The name of the Google Pub/Sub Topic for sensor data"
}

variable "catastrophe_function" {
  default     = "black-box-simulate-catastrophe"
  description = "(Optional) The name of the Google Cloud Function to simulate a catastrophe"
}

variable "catastrophe_topic" {
  default     = "black-box-catastrophe"
  description = "(Optional) The name of the Google Pub/Sub Topic for simulating a catastrophe"
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
resource "google_pubsub_topic" "sensor_data_topic" {
  name = "${var.sensor_data_topic}"
}

# Add the Particle service account as a Pub/Sub publisher to the topic
resource "google_pubsub_topic_iam_binding" "sensor_data_publisher" {
  topic = "${google_pubsub_topic.sensor_data_topic.name}"
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
data "archive_file" "sensor_data_function_bundle" {
  type        = "zip"
  source_dir  = "./sensorFunction"
  output_path = "dist/${var.sensor_data_function}.zip"
}

# Put the Function zip in the bucket
# Adding the sha256 in the name to cache-bust the Function loader
resource "google_storage_bucket_object" "sensor_data_function_object" {
  name   = "${var.sensor_data_function}-${data.archive_file.sensor_data_function_bundle.output_base64sha256}.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "${data.archive_file.sensor_data_function_bundle.output_path}"
}

# Create the Function with the Go 1.11 runtime
# The Function is triggered by the sensor data topic
resource "google_cloudfunctions_function" "sensor_data_function" {
  name                  = "${var.sensor_data_function}"
  available_memory_mb   = 128
  timeout               = 60
  runtime               = "go111"
  source_archive_bucket = "${google_storage_bucket_object.sensor_data_function_object.bucket}"
  source_archive_object = "${google_storage_bucket_object.sensor_data_function_object.name}"

  event_trigger = {
    event_type = "google.pubsub.topic.publish"
    resource   = "${google_pubsub_topic.sensor_data_topic.name}"
  }
  environment_variables = {
    SHEET_ID = "${var.sheet_id}"
  }

  entry_point = "Run"
}

# Create a Pub/Sub Topic for the catastrophe data
resource "google_pubsub_topic" "catastrophe_topic" {
  name = "${var.catastrophe_topic}"
}

# Add the Particle service account as a Pub/Sub publisher to the topic
resource "google_pubsub_topic_iam_binding" "catastrophe_publisher" {
  topic = "${google_pubsub_topic.catastrophe_topic.name}"
  role  = "roles/pubsub.publisher"

  members = [
    "serviceAccount:particle-public@particle-public.iam.gserviceaccount.com",
  ]
}

data "archive_file" "catastrophe_function_bundle" {
  type        = "zip"
  source_dir  = "./catastropheFunction"
  output_path = "dist/${var.catastrophe_function}.zip"
}

# Put the Function zip in the bucket
# Adding the sha256 in the name to cache-bust the Function loader
resource "google_storage_bucket_object" "catastrophe_function_object" {
  name   = "${var.catastrophe_function}-${data.archive_file.catastrophe_function_bundle.output_base64sha256}.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "${data.archive_file.catastrophe_function_bundle.output_path}"
}

# Create the Function with the Go 1.11 runtime
# The Function is triggered by the sensor data topic
resource "google_cloudfunctions_function" "catastrophe_function" {
  name                  = "${var.catastrophe_function}"
  available_memory_mb   = 128
  timeout               = 60
  runtime               = "go111"
  source_archive_bucket = "${google_storage_bucket_object.catastrophe_function_object.bucket}"
  source_archive_object = "${google_storage_bucket_object.catastrophe_function_object.name}"

  event_trigger = {
    event_type = "google.pubsub.topic.publish"
    resource   = "${google_pubsub_topic.catastrophe_topic.name}"
  }
  environment_variables = {
    DEVICE_IDS = "${var.device_ids}"
    PARTICLE_ACCESS_TOKEN = "${var.particle_access_token}"
  }

  entry_point = "Run"
}
