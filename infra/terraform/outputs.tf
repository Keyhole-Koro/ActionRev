output "backend_url" {
  value = module.cloudrun.backend_url
}

output "uploads_bucket_name" {
  value = module.storage.uploads_bucket_name
}

output "sandbox_bucket_name" {
  value = module.storage.sandbox_bucket_name
}

output "bigquery_dataset_id" {
  value = module.bigquery.dataset_id
}

output "artifact_registry_url" {
  value = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.actionrev.repository_id}"
}

output "backend_sa_email" {
  value = module.iam.backend_sa_email
}

output "cloud_tasks_queue" {
  value = module.tasks.queue_name
}
