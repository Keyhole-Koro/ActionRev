resource "google_artifact_registry_repository" "synthify" {
  repository_id = "synthify"
  location      = var.region
  format        = "DOCKER"
  description   = "ActionRev Docker images"
  labels        = var.labels

  cleanup_policies {
    id     = "keep-latest-10"
    action = "KEEP"
    most_recent_versions { keep_count = 10 }
  }

  cleanup_policies {
    id     = "delete-untagged-after-30d"
    action = "DELETE"
    condition {
      tag_state  = "UNTAGGED"
      older_than = "2592000s"
    }
  }
}
