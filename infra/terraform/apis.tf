# 必要な GCP API を有効化する
# 有効化まで数分かかる場合がある

locals {
  required_apis = [
    "run.googleapis.com",               # Cloud Run
    "storage.googleapis.com",           # Cloud Storage
    "bigquery.googleapis.com",          # BigQuery
    "aiplatform.googleapis.com",        # Vertex AI (Gemini)
    "secretmanager.googleapis.com",     # Secret Manager
    "cloudtasks.googleapis.com",        # Cloud Tasks
    "artifactregistry.googleapis.com",  # Artifact Registry
    "cloudbuild.googleapis.com",        # Cloud Build (CI/CD)
    "logging.googleapis.com",           # Cloud Logging
    "monitoring.googleapis.com",        # Cloud Monitoring
    "iam.googleapis.com",               # IAM
    "cloudresourcemanager.googleapis.com",
    "firebase.googleapis.com",          # Firebase
    "identitytoolkit.googleapis.com",   # Firebase Auth (Identity Platform)
  ]
}

resource "google_project_service" "apis" {
  for_each = toset(local.required_apis)

  project                    = var.project_id
  service                    = each.value
  disable_on_destroy         = false
  disable_dependent_services = false
}
