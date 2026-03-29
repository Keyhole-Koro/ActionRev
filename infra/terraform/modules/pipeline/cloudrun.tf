# サンドボックス Cloud Run Job (正規化ツール隔離実行)
resource "google_cloud_run_v2_job" "sandbox" {
  name     = "synthify-sandbox-${var.env}"
  location = var.region
  labels   = var.labels

  template {
    template {
      service_account = var.sandbox_sa_email
      max_retries     = 0      # 冪等性を保証できないためリトライなし
      timeout         = "300s"

      vpc_access { egress = "PRIVATE_RANGES_ONLY" }  # 外部通信禁止

      containers {
        image = "gcr.io/cloudrun/placeholder"  # CI/CD でデプロイ時に上書き
        resources { limits = { cpu = "1", memory = "512Mi" } }
        env { name = "GCS_SANDBOX_BUCKET"; value = google_storage_bucket.sandbox.name }
      }
    }
  }
}
