# ---------------------------------------------------------------------------
# Secret Manager シークレット定義
# 値は terraform apply 後に手動で設定する:
#   gcloud secrets versions add <secret-id> --data-file=-
# ---------------------------------------------------------------------------

resource "google_secret_manager_secret" "stripe_secret_key" {
  secret_id = "stripe-secret-key"

  labels = local.labels

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret" "stripe_webhook_secret" {
  secret_id = "stripe-webhook-secret"

  labels = local.labels

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret" "discord_webhook_url" {
  secret_id = "discord-webhook-url"

  labels = local.labels

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

# Firebase Admin SDK サービスアカウント JSON (バックエンドから Firebase Admin を使う場合)
resource "google_secret_manager_secret" "firebase_admin_credentials" {
  secret_id = "firebase-admin-credentials"

  labels = local.labels

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

# ---------------------------------------------------------------------------
# バックエンド SA にシークレットへのアクセス権を付与
# (IAM でまとめて付与しているが、シークレット個別の ACL が必要な場合はここで設定)
# ---------------------------------------------------------------------------
