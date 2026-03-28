variable "project_id" { type = string }
variable "region"     { type = string }
variable "env"        { type = string }
variable "labels"     { type = map(string) }

variable "uploads_bucket_name" {
  description = "空の場合は <project_id>-uploads-<env> を使用"
  type        = string
  default     = ""
}
