output "queue_name"          { value = google_cloud_tasks_queue.processing.name }
output "sandbox_job_name"    { value = google_cloud_run_v2_job.sandbox.name }
output "sandbox_bucket_name" { value = google_storage_bucket.sandbox.name }
