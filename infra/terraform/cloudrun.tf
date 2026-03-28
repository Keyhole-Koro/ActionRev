# ---------------------------------------------------------------------------
# バックエンド Cloud Run サービス (Connect RPC)
# ---------------------------------------------------------------------------

resource "google_cloud_run_v2_service" "backend" {
  name     = "actionrev-backend-${var.env}"
  location = var.region

  labels = local.labels

  template {
    service_account = google_service_account.backend.email

    scaling {
      min_instance_count = var.backend_min_instances
      max_instance_count = var.backend_max_instances
    }

    containers {
      image = var.backend_image

      ports {
        name           = "h2c"   # HTTP/2 クリアテキスト (Connect RPC)
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        cpu_idle = true
      }

      env {
        name  = "ENV"
        value = var.env
      }
      env {
        name  = "GCP_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "GCS_BUCKET"
        value = google_storage_bucket.uploads.name
      }
      env {
        name  = "BIGQUERY_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "BIGQUERY_DATASET"
        value = google_bigquery_dataset.graph.dataset_id
      }
      env {
        name  = "VERTEX_AI_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "VERTEX_AI_LOCATION"
        value = var.region
      }
      env {
        name  = "CLOUD_TASKS_QUEUE"
        value = google_cloud_tasks_queue.processing.name
      }
      env {
        name  = "CLOUD_TASKS_LOCATION"
        value = var.region
      }
      env {
        name  = "SANDBOX_JOB_NAME"
        value = google_cloud_run_v2_job.sandbox.name
      }

      # Secret Manager からシークレットを注入
      env {
        name = "STRIPE_SECRET_KEY"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.stripe_secret_key.secret_id
            version = "latest"
          }
        }
      }
      env {
        name = "STRIPE_WEBHOOK_SECRET"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.stripe_webhook_secret.secret_id
            version = "latest"
          }
        }
      }
      env {
        name = "DISCORD_WEBHOOK_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.discord_webhook_url.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  depends_on = [
    google_project_service.apis,
    google_project_iam_member.backend,
  ]
}

# Cloud Run サービスを全ユーザーに公開 (Firebase Auth で認証を行うため IAP は使用しない)
resource "google_cloud_run_v2_service_iam_member" "backend_public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# ---------------------------------------------------------------------------
# サンドボックス Cloud Run Job (正規化ツール隔離実行)
# ---------------------------------------------------------------------------

resource "google_cloud_run_v2_job" "sandbox" {
  name     = "actionrev-sandbox-${var.env}"
  location = var.region

  labels = local.labels

  template {
    template {
      service_account = google_service_account.sandbox.email

      # ネットワークを切断してサンドボックス実行
      vpc_access {
        egress = "PRIVATE_RANGES_ONLY"
      }

      max_retries = 0  # 失敗時はリトライしない (冪等性を保証できないため)

      timeout = "300s"  # 5分タイムアウト

      containers {
        image = var.sandbox_image

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }

        env {
          name  = "GCS_SANDBOX_BUCKET"
          value = google_storage_bucket.sandbox.name
        }
      }
    }
  }

  depends_on = [google_project_service.apis]
}
