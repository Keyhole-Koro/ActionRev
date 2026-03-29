# ローカル GCP 開発環境用 Terraform state
# terraform init -backend-config=backends/local.hcl
bucket = "REPLACE_PROJECT_ID-tfstate"
prefix = "actionrev/local"
