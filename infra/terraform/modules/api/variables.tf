variable "project_id"                    { type = string }
variable "region"                        { type = string }
variable "env"                           { type = string }
variable "labels"                        { type = map(string) }
variable "backend_sa_email"              { type = string }
variable "uploads_bucket_name"           { type = string }
variable "bigquery_dataset_id"           { type = string }
variable "cloud_tasks_queue"             { type = string }
variable "sandbox_job_name"              { type = string }
variable "sandbox_bucket_name"           { type = string }
variable "stripe_secret_key_id"          { type = string }
variable "stripe_webhook_secret_id"      { type = string }
variable "discord_webhook_url_id"        { type = string }
variable "firebase_admin_credentials_id" { type = string }
variable "backend_image"                 { type = string }
variable "min_instances"                 { type = number; default = 0 }
variable "max_instances"                 { type = number; default = 10 }
variable "gemini_model"                  { type = string; default = "gemini-2.0-flash-001" }
