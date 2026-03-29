# アップロードファイル保存バケット
resource "google_storage_bucket" "uploads" {
  name                        = "${var.project_id}-uploads"
  location                    = var.region
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  labels                      = var.labels

  lifecycle_rule {
    condition { age = 7; matches_prefix = ["pending/"]; with_state = "ANY" }
    action { type = "Delete" }
  }

  lifecycle_rule {
    condition { age = 90; matches_prefix = ["processed/"] }
    action { type = "SetStorageClass"; storage_class = "NEARLINE" }
  }

  cors {
    origin          = ["*"]
    method          = ["PUT", "GET", "HEAD"]
    response_header = ["Content-Type", "x-goog-meta-*", "ETag"]
    max_age_seconds = 3600
  }
}
