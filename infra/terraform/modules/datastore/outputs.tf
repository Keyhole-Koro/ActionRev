output "dataset_id"          { value = google_bigquery_dataset.graph.dataset_id }
output "uploads_bucket_name" { value = google_storage_bucket.uploads.name }
