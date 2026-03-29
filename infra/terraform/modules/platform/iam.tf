resource "google_service_account" "backend" {
  account_id   = "actionrev-backend"
  display_name = "ActionRev Backend (Cloud Run)"
  project      = var.project_id
}

resource "google_service_account" "sandbox" {
  account_id   = "actionrev-sandbox"
  display_name = "ActionRev Sandbox Runner (Cloud Run Jobs)"
  project      = var.project_id
}

locals {
  backend_roles = [
    "roles/bigquery.dataEditor",
    "roles/bigquery.jobUser",
    "roles/storage.objectAdmin",
    "roles/aiplatform.user",
    "roles/secretmanager.secretAccessor",
    "roles/cloudtasks.enqueuer",
    "roles/run.invoker",
    "roles/iam.serviceAccountTokenCreator",  # 署名付き URL 発行
  ]
}

resource "google_project_iam_member" "backend" {
  for_each = toset(local.backend_roles)
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.backend.email}"
}

data "google_project" "current" { project_id = var.project_id }

resource "google_project_iam_member" "cloudbuild_ar" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${data.google_project.current.number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "cloudbuild_run" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${data.google_project.current.number}@cloudbuild.gserviceaccount.com"
}
