resource "google_cloud_run_v2_service" "backend" {
  name     = "synthify-backend-${var.env}"
  location = var.region
  labels   = var.labels

  template {
    service_account = var.backend_sa_email

    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    containers {
      image = var.backend_image

      ports {
        name           = "h2c"   # HTTP/2 クリアテキスト (Connect RPC)
        container_port = 8080
      }

      resources {
        limits   = { cpu = "1", memory = "512Mi" }
        cpu_idle = true
      }

      env { name = "ENV";                  value = var.env }
      env { name = "GCP_PROJECT_ID";       value = var.project_id }
      env { name = "GCS_BUCKET";           value = var.uploads_bucket_name }
      env { name = "BIGQUERY_PROJECT_ID";  value = var.project_id }
      env { name = "BIGQUERY_DATASET";     value = var.bigquery_dataset_id }
      env { name = "VERTEX_AI_PROJECT_ID"; value = var.project_id }
      env { name = "VERTEX_AI_LOCATION";   value = var.region }
      env { name = "GEMINI_MODEL";         value = var.gemini_model }
      env { name = "CLOUD_TASKS_QUEUE";    value = var.cloud_tasks_queue }
      env { name = "CLOUD_TASKS_LOCATION"; value = var.region }
      env { name = "SANDBOX_JOB_NAME";     value = var.sandbox_job_name }
      env { name = "SANDBOX_BUCKET";       value = var.sandbox_bucket_name }
      env { name = "GEMINI_CACHE_ENABLED"; value = "false" }

      env {
        name = "STRIPE_SECRET_KEY"
        value_source { secret_key_ref { secret = var.stripe_secret_key_id; version = "latest" } }
      }
      env {
        name = "STRIPE_WEBHOOK_SECRET"
        value_source { secret_key_ref { secret = var.stripe_webhook_secret_id; version = "latest" } }
      }
      env {
        name = "DISCORD_WEBHOOK_URL"
        value_source { secret_key_ref { secret = var.discord_webhook_url_id; version = "latest" } }
      }
      env {
        name = "FIREBASE_ADMIN_CREDENTIALS"
        value_source { secret_key_ref { secret = var.firebase_admin_credentials_id; version = "latest" } }
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
