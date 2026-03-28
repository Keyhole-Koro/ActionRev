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
  description = "デプロイ環境 (dev / prod)"
  type        = string
  validation {
    condition     = contains(["dev", "prod"], var.env)
    error_message = "env は dev または prod のみ使用可能。"
  }
}

variable "backend_image" {
  description = "Cloud Run バックエンドの Docker イメージ URI"
  type        = string
  # 例: asia-northeast1-docker.pkg.dev/<project>/actionrev/backend:latest
}

variable "sandbox_image" {
  description = "Cloud Run Jobs サンドボックスの Docker イメージ URI"
  type        = string
}

variable "backend_min_instances" {
  description = "Cloud Run バックエンドの最小インスタンス数 (cold start 対策)"
  type        = number
  default     = 0
}

variable "backend_max_instances" {
  description = "Cloud Run バックエンドの最大インスタンス数"
  type        = number
  default     = 10
}

variable "bigquery_dataset_id" {
  description = "BigQuery データセット ID"
  type        = string
  default     = "graph"
}

variable "gcs_uploads_bucket" {
  description = "アップロードファイル保存バケット名 (空の場合は自動生成)"
  type        = string
  default     = ""
}
