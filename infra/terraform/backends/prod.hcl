# 本番環境用 Terraform state
# terraform init -backend-config=backends/prod.hcl
bucket = "REPLACE_PROJECT_ID-tfstate"
prefix = "actionrev/prod"
