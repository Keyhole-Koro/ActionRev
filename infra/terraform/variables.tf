variable "project_id" {
  description = "GCP プロジェクト ID"
  type        = string
}

variable "region" {
  description = "GCP リージョン"
  type        = string
  default     = "asia-northeast1"
}

variable "env" {
  description = "デプロイ環境"
  type        = string
  validation {
    condition     = contains(["local", "stage", "prod"], var.env)
    error_message = "env は local / stage / prod のみ使用可能。"
  }
}

variable "backend_image" {
  description = "Cloud Run バックエンドの Docker イメージ URI"
  type        = string
}

variable "sandbox_image" {
  description = "Cloud Run Jobs サンドボックスの Docker イメージ URI"
  type        = string
}

variable "backend_min_instances" {
  type    = number
  default = 0
}

variable "backend_max_instances" {
  type    = number
  default = 10
}

variable "bigquery_dataset_id" {
  type    = string
  default = "graph"
}

variable "uploads_bucket_name" {
  description = "空の場合は <project_id>-uploads-<env> を自動生成"
  type        = string
  default     = ""
}

variable "gemini_model" {
  type    = string
  default = "gemini-2.0-flash-001"
}
