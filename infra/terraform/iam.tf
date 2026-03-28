# ---------------------------------------------------------------------------
# サービスアカウント
# ---------------------------------------------------------------------------

resource "google_service_account" "backend" {
  account_id   = "actionrev-backend"
  display_name = "ActionRev Backend (Cloud Run)"
  project      = var.project_id
  depends_on   = [google_project_service.apis]
}

resource "google_service_account" "sandbox" {
  account_id   = "actionrev-sandbox"
  display_name = "ActionRev Sandbox Runner (Cloud Run Jobs)"
  project      = var.project_id
  depends_on   = [google_project_service.apis]
}

# ---------------------------------------------------------------------------
# バックエンド SA の権限
# ---------------------------------------------------------------------------

locals {
  backend_roles = [
    "roles/bigquery.dataEditor",        # BigQuery 読み書き
    "roles/bigquery.jobUser",           # BigQuery ジョブ実行
    "roles/storage.objectAdmin",        # GCS オブジェクト操作 (署名付き URL 発行含む)
    "roles/aiplatform.user",            # Vertex AI / Gemini 呼び出し
    "roles/secretmanager.secretAccessor", # Secret Manager 参照
    "roles/cloudtasks.enqueuer",        # Cloud Tasks エンキュー
    "roles/run.invoker",                # Cloud Run サービス間呼び出し (sandbox 起動)
  ]
}

resource "google_project_iam_member" "backend" {
  for_each = toset(local.backend_roles)

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.backend.email}"
}

# ---------------------------------------------------------------------------
# サンドボックス SA の権限 (最小権限: GCS 読み書きのみ)
# ---------------------------------------------------------------------------

resource "google_storage_bucket_iam_member" "sandbox_read" {
  bucket = google_storage_bucket.sandbox.name
  role   = "roles/storage.objectUser"
  member = "serviceAccount:${google_service_account.sandbox.email}"
}

# ---------------------------------------------------------------------------
# Cloud Build SA に Artifact Registry への push 権限
# ---------------------------------------------------------------------------

resource "google_project_iam_member" "cloudbuild_ar_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"

  depends_on = [google_project_service.apis]
}

resource "google_project_iam_member" "cloudbuild_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"

  depends_on = [google_project_service.apis]
}

data "google_project" "project" {
  project_id = var.project_id
}
