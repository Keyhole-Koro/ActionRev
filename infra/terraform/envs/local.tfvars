# ローカル GCP 開発環境 (個人 GCP プロジェクト)
# Docker Compose での純ローカル開発は infra/scripts/local-setup.sh を使用すること

project_id = "actionrev-dev-YOURNAME"
region     = "asia-northeast1"
env        = "local"

backend_image = "asia-northeast1-docker.pkg.dev/actionrev-dev-YOURNAME/actionrev/backend:latest"

backend_min_instances = 0
backend_max_instances = 3

bigquery_dataset_id = "graph"
gemini_model        = "gemini-2.0-flash-001"
