variable "project_id" {
  description = "GCP プロジェクト ID"
  type        = string
}

variable "region" {
  type    = string
  default = "asia-northeast1"
}

variable "env" {
  type = string
  validation {
    condition     = contains(["local", "stage", "prod"], var.env)
    error_message = "env は local / stage / prod のみ。"
  }
}

variable "backend_image" {
  description = "Cloud Run バックエンドの Docker イメージ URI"
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

variable "gemini_model" {
  type    = string
  default = "gemini-2.0-flash-001"
}
