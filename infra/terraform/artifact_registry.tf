# ---------------------------------------------------------------------------
# Artifact Registry: Docker イメージ管理
# ---------------------------------------------------------------------------

resource "google_artifact_registry_repository" "actionrev" {
  repository_id = "actionrev"
  location      = var.region
  format        = "DOCKER"
  description   = "ActionRev Docker images"

  labels = local.labels

  # 古いイメージの自動クリーンアップ
  cleanup_policies {
    id     = "keep-latest-10"
    action = "KEEP"
    most_recent_versions {
      keep_count = 10
    }
  }

  cleanup_policies {
    id     = "delete-untagged-after-30d"
    action = "DELETE"
    condition {
      tag_state  = "UNTAGGED"
      older_than = "2592000s"  # 30日
    }
  }

  depends_on = [google_project_service.apis]
}
