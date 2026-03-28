output "backend_url" {
  description = "Cloud Run バックエンドの URL"
  value       = google_cloud_run_v2_service.backend.uri
}

output "uploads_bucket_name" {
  description = "アップロードファイル用 GCS バケット名"
  value       = google_storage_bucket.uploads.name
}

output "sandbox_bucket_name" {
  description = "サンドボックス作業用 GCS バケット名"
  value       = google_storage_bucket.sandbox.name
}

output "bigquery_dataset" {
  description = "BigQuery データセット ID"
  value       = google_bigquery_dataset.graph.dataset_id
}

output "artifact_registry_repository" {
  description = "Artifact Registry リポジトリ URL"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.actionrev.repository_id}"
}

output "backend_service_account_email" {
  description = "Cloud Run バックエンドのサービスアカウント"
  value       = google_service_account.backend.email
}

output "cloud_tasks_queue_name" {
  description = "Cloud Tasks キュー名"
  value       = google_cloud_tasks_queue.processing.name
}
