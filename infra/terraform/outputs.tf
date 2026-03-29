output "backend_url"             { value = module.api.backend_url }
output "uploads_bucket_name"     { value = module.datastore.uploads_bucket_name }
output "sandbox_bucket_name"     { value = module.pipeline.sandbox_bucket_name }
output "bigquery_dataset_id"     { value = module.datastore.dataset_id }
output "artifact_registry_url"   { value = module.platform.artifact_registry_url }
output "backend_sa_email"        { value = module.platform.backend_sa_email }
output "cloud_tasks_queue"       { value = module.pipeline.queue_name }
