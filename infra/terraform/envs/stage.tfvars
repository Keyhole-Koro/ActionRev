# ステージング環境

project_id = "actionrev-stage"
region     = "asia-northeast1"
env        = "stage"

backend_image = "asia-northeast1-docker.pkg.dev/actionrev-stage/actionrev/backend:latest"

backend_min_instances = 0    # コスト抑制のためコールドスタート許容
backend_max_instances = 5

bigquery_dataset_id = "graph"
gemini_model        = "gemini-2.0-flash-001"
