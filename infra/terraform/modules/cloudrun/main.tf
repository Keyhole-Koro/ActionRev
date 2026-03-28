# ---------------------------------------------------------------------------
# バックエンド Cloud Run サービス (Connect RPC / HTTP2)
# ---------------------------------------------------------------------------

resource "google_cloud_run_v2_service" "backend" {
  name     = "actionrev-backend-${var.env}"
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

      # --- Static env vars ---
      env { name = "ENV";                   value = var.env }
      env { name = "GCP_PROJECT_ID";        value = var.project_id }
      env { name = "GCS_BUCKET";            value = var.uploads_bucket_name }
      env { name = "BIGQUERY_PROJECT_ID";   value = var.project_id }
      env { name = "BIGQUERY_DATASET";      value = var.bigquery_dataset_id }
      env { name = "VERTEX_AI_PROJECT_ID";  value = var.project_id }
      env { name = "VERTEX_AI_LOCATION";    value = var.region }
      env { name = "GEMINI_MODEL";          value = var.gemini_model }
      env { name = "CLOUD_TASKS_QUEUE";     value = var.cloud_tasks_queue }
      env { name = "CLOUD_TASKS_LOCATION";  value = var.region }
      env { name = "SANDBOX_JOB_NAME";      value = google_cloud_run_v2_job.sandbox.name }
      env { name = "SANDBOX_BUCKET";        value = var.sandbox_bucket_name }

      # --- Gemini response cache (本番では無効) ---
      env { name = "GEMINI_CACHE_ENABLED";  value = "false" }

      # --- Secrets from Secret Manager ---
      env {
        name = "STRIPE_SECRET_KEY"
        value_source {
          secret_key_ref { secret = var.stripe_secret_key_id; version = "latest" }
        }
      }
      env {
        name = "STRIPE_WEBHOOK_SECRET"
        value_source {
          secret_key_ref { secret = var.stripe_webhook_secret_id; version = "latest" }
        }
      }
      env {
        name = "DISCORD_WEBHOOK_URL"
        value_source {
          secret_key_ref { secret = var.discord_webhook_url_id; version = "latest" }
        }
      }
      env {
        name = "FIREBASE_ADMIN_CREDENTIALS"
        value_source {
          secret_key_ref { secret = var.firebase_admin_credentials_id; version = "latest" }
        }
      }
    }
  }
}

# 全ユーザーに公開 (Firebase Auth で認可するためフロントエンドからの直接アクセスを許可)
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
  labels   = var.labels

  template {
    template {
      service_account = var.sandbox_sa_email
      max_retries     = 0       # 冪等性を保証できないためリトライなし
      timeout         = "300s"

      # ネットワーク無効 (外部通信禁止)
      vpc_access { egress = "PRIVATE_RANGES_ONLY" }

      containers {
        image = var.sandbox_image
        resources { limits = { cpu = "1", memory = "512Mi" } }
        env { name = "GCS_SANDBOX_BUCKET"; value = var.sandbox_bucket_name }
      }
    }
  }
}
