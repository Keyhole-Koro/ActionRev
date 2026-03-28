# サンドボックス作業バケット (正規化ツールの入出力)
resource "google_storage_bucket" "sandbox" {
  name                        = "${var.project_id}-sandbox-${var.env}"
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  labels                      = var.labels

  # 作業ファイルは 1 日で自動削除
  lifecycle_rule {
    condition { age = 1 }
    action { type = "Delete" }
  }
}

# sandbox SA にバケットへのアクセス権を付与
# (platform モジュールで SA を作成後、pipeline モジュールでバケット IAM を設定する)
resource "google_storage_bucket_iam_member" "sandbox_sa" {
  bucket = google_storage_bucket.sandbox.name
  role   = "roles/storage.objectUser"
  member = "serviceAccount:${var.sandbox_sa_email}"
}
