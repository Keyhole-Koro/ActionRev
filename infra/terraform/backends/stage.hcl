# ステージング環境用 Terraform state
# terraform init -backend-config=backends/stage.hcl
bucket = "REPLACE_PROJECT_ID-tfstate"
prefix = "actionrev/stage"
