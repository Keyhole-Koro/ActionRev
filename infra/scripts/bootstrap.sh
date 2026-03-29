#!/usr/bin/env bash
# bootstrap.sh - GCP プロジェクトの初期セットアップ
# Terraform を実行する前に一度だけ実行する
#
# Usage:
#   ./infra/scripts/bootstrap.sh <project-id> <env>
#   ./infra/scripts/bootstrap.sh my-project-id stage
#
set -euo pipefail

PROJECT_ID="${1:-}"
ENV="${2:-stage}"

if [[ -z "$PROJECT_ID" ]]; then
  echo "Usage: $0 <project-id> <env>"
  echo "  env: local | stage | prod"
  exit 1
fi

REGION="asia-northeast1"
TFSTATE_BUCKET="${PROJECT_ID}-tfstate"
TF_SA="terraform@${PROJECT_ID}.iam.gserviceaccount.com"

echo "=== ActionRev Bootstrap: project=${PROJECT_ID} env=${ENV} ==="

# ---------------------------------------------------------------------------
# 1. 最低限必要な API を有効化 (Terraform 実行前に必要)
# ---------------------------------------------------------------------------
echo ">> Enabling core APIs..."
gcloud services enable \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com \
  storage.googleapis.com \
  --project="$PROJECT_ID"

# ---------------------------------------------------------------------------
# 2. Terraform state 用 GCS バケットを作成
# ---------------------------------------------------------------------------
echo ">> Creating Terraform state bucket: gs://${TFSTATE_BUCKET}"
if ! gsutil ls -p "$PROJECT_ID" "gs://${TFSTATE_BUCKET}" &>/dev/null; then
  gsutil mb -p "$PROJECT_ID" -l "$REGION" "gs://${TFSTATE_BUCKET}"
  gsutil versioning set on "gs://${TFSTATE_BUCKET}"
  gsutil lifecycle set - "gs://${TFSTATE_BUCKET}" <<'LIFECYCLE'
{
  "rule": [
    {
      "action": {"type": "Delete"},
      "condition": {"numNewerVersions": 5, "isLive": false}
    }
  ]
}
LIFECYCLE
  echo "   Created: gs://${TFSTATE_BUCKET}"
else
  echo "   Already exists: gs://${TFSTATE_BUCKET}"
fi

# ---------------------------------------------------------------------------
# 3. Terraform 実行用サービスアカウントを作成
# ---------------------------------------------------------------------------
echo ">> Creating Terraform service account: ${TF_SA}"
if ! gcloud iam service-accounts describe "$TF_SA" --project="$PROJECT_ID" &>/dev/null; then
  gcloud iam service-accounts create terraform \
    --display-name="Terraform Runner" \
    --project="$PROJECT_ID"
  echo "   Created: ${TF_SA}"
else
  echo "   Already exists: ${TF_SA}"
fi

# Terraform SA にプロジェクト編集権限を付与
# 本番では roles/owner ではなく必要最小限のロールに絞ることを推奨
for ROLE in \
  "roles/editor" \
  "roles/iam.securityAdmin" \
  "roles/resourcemanager.projectIamAdmin"; do
  gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${TF_SA}" \
    --role="$ROLE" \
    --quiet
done

# Terraform SA に tfstate バケットへのアクセス権を付与
gsutil iam ch "serviceAccount:${TF_SA}:roles/storage.objectAdmin" "gs://${TFSTATE_BUCKET}"

# ---------------------------------------------------------------------------
# 4. Terraform SA のキーを生成 (CI/CD 用。不要なら --key-file を省略)
# ---------------------------------------------------------------------------
KEY_FILE="infra/secrets/terraform-sa-${ENV}.json"
mkdir -p "$(dirname "$KEY_FILE")"

if [[ ! -f "$KEY_FILE" ]]; then
  echo ">> Generating SA key: ${KEY_FILE}"
  gcloud iam service-accounts keys create "$KEY_FILE" \
    --iam-account="$TF_SA" \
    --project="$PROJECT_ID"
  echo "   IMPORTANT: ${KEY_FILE} を安全な場所に保管し、git に commit しないこと"
else
  echo "   Key already exists: ${KEY_FILE}"
fi

# ---------------------------------------------------------------------------
# 5. backends/ の project_id を置換
# ---------------------------------------------------------------------------
BACKEND_FILE="infra/terraform/backends/${ENV}.hcl"
if [[ -f "$BACKEND_FILE" ]]; then
  sed -i "s/REPLACE_PROJECT_ID/${PROJECT_ID}/g" "$BACKEND_FILE"
  echo ">> Updated backend config: ${BACKEND_FILE}"
fi

echo ""
echo "=== Bootstrap complete! ==="
echo ""
echo "Next steps:"
echo "  1. cd infra/terraform"
echo "  2. terraform init -backend-config=backends/${ENV}.hcl"
echo "  3. terraform apply -var-file=envs/${ENV}.tfvars"
echo "  4. After apply: ./infra/scripts/set-secrets.sh ${PROJECT_ID} ${ENV}"
