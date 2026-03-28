#!/usr/bin/env bash
# set-secrets.sh - Secret Manager にシークレット値を設定する
# Terraform apply 後に一度実行する
#
# Usage:
#   ./infra/scripts/set-secrets.sh <project-id>
#
set -euo pipefail

PROJECT_ID="${1:-}"

if [[ -z "$PROJECT_ID" ]]; then
  echo "Usage: $0 <project-id>"
  exit 1
fi

echo "=== Setting secrets for project: ${PROJECT_ID} ==="
echo "各シークレットの値を入力してください (空Enter でスキップ)"
echo ""

set_secret() {
  local SECRET_ID="$1"
  local PROMPT="$2"
  local DEFAULT="${3:-}"

  if [[ -n "$DEFAULT" ]]; then
    read -rsp "${PROMPT} [${DEFAULT}]: " VALUE
    VALUE="${VALUE:-$DEFAULT}"
  else
    read -rsp "${PROMPT}: " VALUE
  fi
  echo ""

  if [[ -z "$VALUE" ]]; then
    echo "  Skipped: ${SECRET_ID}"
    return
  fi

  echo -n "$VALUE" | gcloud secrets versions add "$SECRET_ID" \
    --data-file=- \
    --project="$PROJECT_ID"
  echo "  Set: ${SECRET_ID}"
}

set_secret "stripe-secret-key"        "Stripe Secret Key (sk_live_...)"
set_secret "stripe-webhook-secret"    "Stripe Webhook Secret (whsec_...)"
set_secret "discord-webhook-url"      "Discord Webhook URL"
set_secret "firebase-admin-credentials" "Firebase Admin SDK JSON (paste single line)"

echo ""
echo "=== Secrets configured! ==="
echo ""
echo "Verify with:"
echo "  gcloud secrets list --project=${PROJECT_ID}"
