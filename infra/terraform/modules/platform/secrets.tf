resource "google_secret_manager_secret" "stripe_secret_key" {
  secret_id = "stripe-secret-key"
  labels    = var.labels
  replication { auto {} }
}

resource "google_secret_manager_secret" "stripe_webhook_secret" {
  secret_id = "stripe-webhook-secret"
  labels    = var.labels
  replication { auto {} }
}

resource "google_secret_manager_secret" "discord_webhook_url" {
  secret_id = "discord-webhook-url"
  labels    = var.labels
  replication { auto {} }
}

resource "google_secret_manager_secret" "firebase_admin_credentials" {
  secret_id = "firebase-admin-credentials"
  labels    = var.labels
  replication { auto {} }
}
