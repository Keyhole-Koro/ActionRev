# 本番環境

project_id = "actionrev-prod"
region     = "asia-northeast1"
env        = "prod"

backend_image = "asia-northeast1-docker.pkg.dev/actionrev-prod/actionrev/backend:latest"

backend_min_instances = 0
backend_max_instances = 20

bigquery_dataset_id = "graph"
gemini_model        = "gemini-2.0-flash-001"
