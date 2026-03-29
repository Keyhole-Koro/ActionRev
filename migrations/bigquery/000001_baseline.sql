-- 000001_baseline.sql
-- description: Baseline - tables created by Terraform on initial deploy
--
-- このファイルは Terraform で作成済みのテーブルを migration 管理の起点として記録する。
-- SQL 自体はすべて IF NOT EXISTS のため冪等だが、
-- Terraform apply 済みの環境では実質 no-op になる。

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.users` (
  user_id       STRING    NOT NULL,
  email         STRING,
  display_name  STRING,
  created_at    TIMESTAMP,
  last_login_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.workspaces` (
  workspace_id           STRING    NOT NULL,
  name                   STRING,
  owner_id               STRING,
  plan                   STRING,
  stripe_customer_id     STRING,
  stripe_subscription_id STRING,
  storage_used_bytes     INT64,
  created_at             TIMESTAMP,
  updated_at             TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.workspace_members` (
  workspace_id STRING    NOT NULL,
  user_id      STRING    NOT NULL,
  role         STRING,
  invited_at   TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.documents` (
  document_id      STRING    NOT NULL,
  workspace_id     STRING,
  uploaded_by      STRING,
  filename         STRING,
  gcs_uri          STRING,
  mime_type        STRING,
  file_size        INT64,
  status           STRING,
  extraction_depth STRING,
  created_at       TIMESTAMP,
  updated_at       TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.document_chunks` (
  document_id          STRING NOT NULL,
  chunk_id             STRING NOT NULL,
  chunk_index          INT64,
  text                 STRING,
  source_filename      STRING,
  source_page          INT64,
  source_offset_start  INT64,
  source_offset_end    INT64
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.nodes` (
  document_id     STRING    NOT NULL,
  node_id         STRING    NOT NULL,
  label           STRING,
  level           INT64,
  category        STRING,
  entity_type     STRING,
  description     STRING,
  summary_html    STRING,
  source_chunk_id STRING,
  confidence      FLOAT64,
  created_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.edges` (
  document_id     STRING    NOT NULL,
  edge_id         STRING    NOT NULL,
  source_node_id  STRING,
  target_node_id  STRING,
  edge_type       STRING,
  description     STRING,
  weight          FLOAT64,
  source_chunk_id STRING,
  created_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.normalization_tools` (
  tool_id           STRING    NOT NULL,
  name              STRING,
  version           STRING,
  description       STRING,
  problem_pattern   STRING,
  approval_status   STRING,
  approved_by       STRING,
  llm_review_score  FLOAT64,
  llm_review_reason STRING,
  created_by        STRING,
  created_at        TIMESTAMP,
  updated_at        TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.normalization_tool_runs` (
  run_id        STRING    NOT NULL,
  tool_id       STRING,
  document_id   STRING,
  run_type      STRING,
  status        STRING,
  diff_summary  STRING,
  error_message STRING,
  started_at    TIMESTAMP,
  completed_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.plans` (
  plan                      STRING NOT NULL,
  storage_quota_bytes       INT64,
  max_file_size_bytes       INT64,
  max_uploads_per_day       INT64,
  max_members               INT64,
  allowed_extraction_depths STRING
);
