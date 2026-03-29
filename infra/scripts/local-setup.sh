#!/usr/bin/env bash
# local-setup.sh - ローカル開発環境のセットアップ
#
# Usage:
#   ./infra/scripts/local-setup.sh
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
INFRA_DIR="${REPO_ROOT}/infra"

echo "=== ActionRev Local Dev Setup ==="

# ---------------------------------------------------------------------------
# 1. .env.local のコピー
# ---------------------------------------------------------------------------
ENV_EXAMPLE="${INFRA_DIR}/.env.local.example"
ENV_LOCAL="${INFRA_DIR}/.env.local"

if [[ ! -f "$ENV_LOCAL" ]]; then
  cp "$ENV_EXAMPLE" "$ENV_LOCAL"
  echo ">> Created ${ENV_LOCAL} - 必要な値を編集してください (Stripe / Discord など)"
else
  echo ">> ${ENV_LOCAL} は既に存在します (スキップ)"
fi

# ---------------------------------------------------------------------------
# 2. Gemini キャッシュディレクトリを作成
# ---------------------------------------------------------------------------
CACHE_DIR="${REPO_ROOT}/.gemini-cache"
mkdir -p "$CACHE_DIR"
echo ">> Gemini cache dir: ${CACHE_DIR}"

# .gitignore に追加 (未登録の場合)
GITIGNORE="${REPO_ROOT}/.gitignore"
if ! grep -q "^\.gemini-cache" "$GITIGNORE" 2>/dev/null; then
  echo "" >> "$GITIGNORE"
  echo "# Gemini response cache (local dev)" >> "$GITIGNORE"
  echo ".gemini-cache/" >> "$GITIGNORE"
  echo "infra/.env.local" >> "$GITIGNORE"
  echo "infra/secrets/" >> "$GITIGNORE"
  echo ">> Updated .gitignore"
fi

# ---------------------------------------------------------------------------
# 3. Docker Compose を起動
# ---------------------------------------------------------------------------
echo ">> Starting Docker Compose..."
cd "$INFRA_DIR"
docker compose pull --quiet
docker compose up -d

echo ">> Waiting for services to be ready..."
sleep 5

# BigQuery emulator に schema を初期化
echo ">> Initializing BigQuery schema..."
MAX_RETRIES=20
for i in $(seq 1 $MAX_RETRIES); do
  if curl -sf "http://localhost:9050/bigquery/v2/projects/actionrev-local/datasets" > /dev/null 2>&1; then
    break
  fi
  echo "   Waiting for BigQuery emulator... (${i}/${MAX_RETRIES})"
  sleep 3
done

# BigQuery データセットを作成 (emulator はコマンド引数で作るが念のため)
curl -s -X POST "http://localhost:9050/bigquery/v2/projects/actionrev-local/datasets" \
  -H "Content-Type: application/json" \
  -d '{"datasetReference":{"projectId":"actionrev-local","datasetId":"graph"}}' \
  > /dev/null || true

# ---------------------------------------------------------------------------
# 4. fake-gcs にローカルバケットを作成
# ---------------------------------------------------------------------------
echo ">> Creating local GCS buckets..."
GCS_BASE="http://localhost:4443/storage/v1/b"

create_bucket() {
  local BUCKET_NAME="$1"
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "${GCS_BASE}?project=actionrev-local" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"${BUCKET_NAME}\"}")
  if [[ "$STATUS" == "200" ]] || [[ "$STATUS" == "409" ]]; then
    echo "   Bucket: ${BUCKET_NAME} (${STATUS})"
  else
    echo "   WARN: Failed to create bucket ${BUCKET_NAME} (HTTP ${STATUS})"
  fi
}

for BUCKET in "actionrev-uploads" "actionrev-sandbox"; do
  create_bucket "$BUCKET"
done

# ---------------------------------------------------------------------------
# 5. 完了メッセージ
# ---------------------------------------------------------------------------
echo ""
echo "=== Local setup complete! ==="
echo ""
echo "Services:"
echo "  Backend API:       http://localhost:8080"
echo "  BigQuery emulator: http://localhost:9050"
echo "  GCS emulator:      http://localhost:4443"
echo "  Firebase Auth:     http://localhost:9099"
echo "  Firebase UI:       http://localhost:4000"
echo ""
echo "Gemini response cache:"
echo "  Directory: ${CACHE_DIR}"
echo "  Toggle:    GEMINI_CACHE_ENABLED=true in .env.local"
echo ""
echo "To stop: docker compose -f infra/docker-compose.yml down"
