locals {
  uploads_bucket_name = var.gcs_uploads_bucket != "" ? var.gcs_uploads_bucket : "${var.project_id}-uploads-${var.env}"
}

# ---------------------------------------------------------------------------
# アップロードファイル保存バケット
# ---------------------------------------------------------------------------

resource "google_storage_bucket" "uploads" {
  name                        = local.uploads_bucket_name
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"  # 公開アクセス禁止

  labels = local.labels

  # ライフサイクル: 未処理のアップロードを 7 日で削除 (アップロードハック対策)
  lifecycle_rule {
    condition {
      age                   = 7
      matches_prefix        = ["pending/"]
      with_state            = "ANY"
    }
    action {
      type = "Delete"
    }
  }

  # ライフサイクル: 処理済みファイルを Nearline に移行 (コスト最適化)
  lifecycle_rule {
    condition {
      age            = 90
      matches_prefix = ["processed/"]
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }

  versioning {
    enabled = false
  }

  cors {
    origin          = ["*"]
    method          = ["PUT", "GET"]
    response_header = ["Content-Type", "x-goog-meta-*"]
    max_age_seconds = 3600
  }

  depends_on = [google_project_service.apis]
}

# ---------------------------------------------------------------------------
# サンドボックス作業バケット (正規化ツールの入出力)
# ---------------------------------------------------------------------------

resource "google_storage_bucket" "sandbox" {
  name                        = "${var.project_id}-sandbox-${var.env}"
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  labels = local.labels

  # サンドボックス作業ファイルは 1 日で自動削除
  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }

  depends_on = [google_project_service.apis]
}
