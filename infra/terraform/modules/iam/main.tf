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
    "roles/storage.objectAdmin",       # 署名付き URL 発行 (iam.serviceAccounts.signBlob)
    "roles/aiplatform.user",           # Vertex AI / Gemini
    "roles/secretmanager.secretAccessor",
    "roles/cloudtasks.enqueuer",
    "roles/run.invoker",               # Cloud Run Job 起動
    "roles/iam.serviceAccountTokenCreator", # 署名付き URL 発行に必要
  ]
}

resource "google_project_iam_member" "backend" {
  for_each = toset(local.backend_roles)
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.backend.email}"
}

# sandbox SA はサンドボックスバケットのみアクセス可
resource "google_storage_bucket_iam_member" "sandbox_bucket" {
  bucket = var.sandbox_bucket_name
  role   = "roles/storage.objectUser"
  member = "serviceAccount:${google_service_account.sandbox.email}"
}

# Cloud Build に Artifact Registry write + Cloud Run デプロイ権限
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
