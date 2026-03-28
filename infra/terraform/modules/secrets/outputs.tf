output "stripe_secret_key_id"          { value = google_secret_manager_secret.stripe_secret_key.secret_id }
output "stripe_webhook_secret_id"      { value = google_secret_manager_secret.stripe_webhook_secret.secret_id }
output "discord_webhook_url_id"        { value = google_secret_manager_secret.discord_webhook_url.secret_id }
output "firebase_admin_credentials_id" { value = google_secret_manager_secret.firebase_admin_credentials.secret_id }
