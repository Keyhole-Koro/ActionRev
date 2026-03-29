# GitHub Secrets 設定ガイド

## Workload Identity Federation (WIF) セットアップ

キーファイル不要の GCP 認証。`infra/scripts/bootstrap.sh` 実行後に以下を行う。

```bash
PROJECT_ID_STAGE="actionrev-stage"
PROJECT_ID_PROD="actionrev-prod"
REPO="keyhole-koro/actionrev"

# WIF Pool 作成 (stage)
gcloud iam workload-identity-pools create github \
  --location=global --project=$PROJECT_ID_STAGE

# Provider 作成 (stage)
gcloud iam workload-identity-pools providers create-oidc github-provider \
  --workload-identity-pool=github \
  --location=global \
  --project=$PROJECT_ID_STAGE \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository=='${REPO}'"

# SA に WIF バインディング (stage)
gcloud iam service-accounts add-iam-policy-binding \
  "actionrev-backend@${PROJECT_ID_STAGE}.iam.gserviceaccount.com" \
  --project=$PROJECT_ID_STAGE \
  --role=roles/iam.workloadIdentityUser \
  --member="principalSet://iam.googleapis.com/projects/$(gcloud projects describe $PROJECT_ID_STAGE --format='value(projectNumber)')/locations/global/workloadIdentityPools/github/attribute.repository/${REPO}"
```

prod も同様に実施。

---

## Repository Secrets (Settings → Secrets → Actions)

| Secret 名 | 値 | 用途 |
|---|---|---|
| `WIF_PROVIDER` | `projects/<number>/locations/global/workloadIdentityPools/github/providers/github-provider` | WIF プロバイダ URI |
| `WIF_SERVICE_ACCOUNT_STAGE` | `actionrev-backend@actionrev-stage.iam.gserviceaccount.com` | stage 用 SA |
| `WIF_SERVICE_ACCOUNT_PROD` | `actionrev-backend@actionrev-prod.iam.gserviceaccount.com` | prod 用 SA |
| `GCP_PROJECT_ID_STAGE` | `actionrev-stage` | stage プロジェクト ID |
| `GCP_PROJECT_ID_PROD` | `actionrev-prod` | prod プロジェクト ID |
| `DISCORD_WEBHOOK_URL` | Discord Incoming Webhook URL | デプロイ通知 |

## GitHub Environments

`Settings → Environments` で以下を作成する。

| Environment | 用途 | Required reviewers |
|---|---|---|
| `stage` | stage デプロイ | なし |
| `prod-plan` | prod の terraform plan 確認 | なし |
| `prod` | prod デプロイ (承認必須) | 1名以上 |
