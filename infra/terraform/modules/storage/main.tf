locals {
  uploads_bucket_name = var.uploads_bucket_name != "" ? var.uploads_bucket_name : "${var.project_id}-uploads-${var.env}"
}

resource "google_storage_bucket" "uploads" {
  name                        = local.uploads_bucket_name
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  labels                      = var.labels

  # pending/ 以下の未処理ファイルを 7 日で削除 (ストレージハック対策)
  lifecycle_rule {
    condition {
      age            = 7
      matches_prefix = ["pending/"]
      with_state     = "ANY"
    }
    action { type = "Delete" }
  }

  # 処理済みファイルを 90 日後に Nearline へ移行 (コスト最適化)
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

  # フロントエンドからの直接 PUT (署名付き URL) に必要な CORS
  cors {
    origin          = ["*"]
    method          = ["PUT", "GET", "HEAD"]
    response_header = ["Content-Type", "x-goog-meta-*", "ETag"]
    max_age_seconds = 3600
  }
}

resource "google_storage_bucket" "sandbox" {
  name                        = "${var.project_id}-sandbox-${var.env}"
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  labels                      = var.labels

  # サンドボックス作業ファイルは 1 日で自動削除
  lifecycle_rule {
    condition { age = 1 }
    action { type = "Delete" }
  }
}
