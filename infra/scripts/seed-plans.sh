#!/usr/bin/env bash
# seed-plans.sh - BigQuery plans テーブルにデフォルトプランをシードする
#
# Usage:
#   ./infra/scripts/seed-plans.sh <project-id> [dataset-id]
#
#   # ローカル (emulator)
#   BIGQUERY_EMULATOR_HOST=localhost:9050 ./infra/scripts/seed-plans.sh actionrev-local graph
#
set -euo pipefail

PROJECT_ID="${1:-}"
DATASET_ID="${2:-graph}"

if [[ -z "$PROJECT_ID" ]]; then
  echo "Usage: $0 <project-id> [dataset-id]"
  exit 1
fi

echo "=== Seeding plans table: ${PROJECT_ID}.${DATASET_ID}.plans ==="

# bq コマンド用フラグ (エミュレータ接続時は環境変数 BIGQUERY_EMULATOR_HOST を設定)
BQ_FLAGS="--project_id=${PROJECT_ID} --dataset_id=${DATASET_ID} --noflag_file"

# plans テーブルのデータを定義
# storage_quota_bytes: free=1GB, pro=50GB
# max_file_size_bytes: free=50MB, pro=500MB
# max_uploads_per_day: free=10, pro=200
# max_members: free=3, pro=20 (0=unlimited)
read -r -d '' PLANS_JSON <<'JSON' || true
{"plan":"free","storage_quota_bytes":1073741824,"max_file_size_bytes":52428800,"max_uploads_per_day":10,"max_members":3,"allowed_extraction_depths":"summary"}
{"plan":"pro","storage_quota_bytes":53687091200,"max_file_size_bytes":524288000,"max_uploads_per_day":200,"max_members":20,"allowed_extraction_depths":"full,summary"}
JSON

echo "$PLANS_JSON" | bq insert $BQ_FLAGS "${DATASET_ID}.plans"

echo ""
echo "Plans seeded:"
bq query $BQ_FLAGS --nouse_legacy_sql \
  "SELECT plan, storage_quota_bytes, max_file_size_bytes, max_uploads_per_day, max_members, allowed_extraction_depths FROM \`${PROJECT_ID}.${DATASET_ID}.plans\` ORDER BY plan"
