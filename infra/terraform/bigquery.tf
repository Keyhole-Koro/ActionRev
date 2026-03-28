# ---------------------------------------------------------------------------
# データセット
# ---------------------------------------------------------------------------

resource "google_bigquery_dataset" "graph" {
  dataset_id                  = var.bigquery_dataset_id
  friendly_name               = "ActionRev Graph"
  location                    = var.region
  delete_contents_on_destroy  = false

  labels = local.labels

  depends_on = [google_project_service.apis]
}

# ---------------------------------------------------------------------------
# テーブル定義
# ---------------------------------------------------------------------------

resource "google_bigquery_table" "users" {
  dataset_id          = google_bigquery_dataset.graph.dataset_id
  table_id            = "users"
  deletion_protection = true

  schema = jsonencode([
    { name = "user_id",       type = "STRING",    mode = "REQUIRED", description = "Firebase Auth UID" },
    { name = "email",         type = "STRING",    mode = "NULLABLE" },
    { name = "display_name",  type = "STRING",    mode = "NULLABLE" },
    { name = "created_at",    type = "TIMESTAMP", mode = "NULLABLE", description = "初回ログイン日時" },
    { name = "last_login_at", type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "MONTH"
    field = "created_at"
  }
}

resource "google_bigquery_table" "workspaces" {
  dataset_id          = google_bigquery_dataset.graph.dataset_id
  table_id            = "workspaces"
  deletion_protection = true

  schema = jsonencode([
    { name = "workspace_id",            type = "STRING",    mode = "REQUIRED" },
    { name = "name",                    type = "STRING",    mode = "NULLABLE" },
    { name = "owner_id",                type = "STRING",    mode = "NULLABLE", description = "Firebase Auth UID" },
    { name = "plan",                    type = "STRING",    mode = "NULLABLE", description = "free / pro" },
    { name = "stripe_customer_id",      type = "STRING",    mode = "NULLABLE" },
    { name = "stripe_subscription_id",  type = "STRING",    mode = "NULLABLE" },
    { name = "storage_used_bytes",      type = "INT64",     mode = "NULLABLE" },
    { name = "created_at",              type = "TIMESTAMP", mode = "NULLABLE" },
    { name = "updated_at",              type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "MONTH"
    field = "created_at"
  }
}

resource "google_bigquery_table" "workspace_members" {
  dataset_id = google_bigquery_dataset.graph.dataset_id
  table_id   = "workspace_members"

  schema = jsonencode([
    { name = "workspace_id", type = "STRING",    mode = "REQUIRED" },
    { name = "user_id",      type = "STRING",    mode = "REQUIRED" },
    { name = "role",         type = "STRING",    mode = "NULLABLE", description = "editor / viewer / dev" },
    { name = "invited_at",   type = "TIMESTAMP", mode = "NULLABLE" },
  ])
}

resource "google_bigquery_table" "documents" {
  dataset_id          = google_bigquery_dataset.graph.dataset_id
  table_id            = "documents"
  deletion_protection = true

  schema = jsonencode([
    { name = "document_id",        type = "STRING",    mode = "REQUIRED" },
    { name = "workspace_id",       type = "STRING",    mode = "NULLABLE" },
    { name = "uploaded_by",        type = "STRING",    mode = "NULLABLE", description = "Firebase Auth UID" },
    { name = "filename",           type = "STRING",    mode = "NULLABLE" },
    { name = "gcs_uri",            type = "STRING",    mode = "NULLABLE" },
    { name = "mime_type",          type = "STRING",    mode = "NULLABLE" },
    { name = "file_size",          type = "INT64",     mode = "NULLABLE" },
    { name = "status",             type = "STRING",    mode = "NULLABLE", description = "uploaded / pending_normalization / processing / completed / failed" },
    { name = "extraction_depth",   type = "STRING",    mode = "NULLABLE", description = "full / summary" },
    { name = "created_at",         type = "TIMESTAMP", mode = "NULLABLE" },
    { name = "updated_at",         type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "DAY"
    field = "created_at"
  }
}

resource "google_bigquery_table" "document_chunks" {
  dataset_id = google_bigquery_dataset.graph.dataset_id
  table_id   = "document_chunks"

  schema = jsonencode([
    { name = "document_id",        type = "STRING", mode = "REQUIRED" },
    { name = "chunk_id",           type = "STRING", mode = "REQUIRED" },
    { name = "chunk_index",        type = "INT64",  mode = "NULLABLE" },
    { name = "text",               type = "STRING", mode = "NULLABLE" },
    { name = "source_filename",    type = "STRING", mode = "NULLABLE", description = "zip 展開後ファイル名。単ファイルの場合は filename と同値" },
    { name = "source_page",        type = "INT64",  mode = "NULLABLE" },
    { name = "source_offset_start",type = "INT64",  mode = "NULLABLE" },
    { name = "source_offset_end",  type = "INT64",  mode = "NULLABLE" },
  ])

  range_partitioning {
    field = "chunk_index"
    range {
      start    = 0
      end      = 100000
      interval = 1000
    }
  }
}

resource "google_bigquery_table" "nodes" {
  dataset_id          = google_bigquery_dataset.graph.dataset_id
  table_id            = "nodes"
  deletion_protection = true

  schema = jsonencode([
    { name = "document_id",    type = "STRING",  mode = "REQUIRED" },
    { name = "node_id",        type = "STRING",  mode = "REQUIRED" },
    { name = "label",          type = "STRING",  mode = "NULLABLE" },
    { name = "level",          type = "INT64",   mode = "NULLABLE", description = "0=ドメイン / 1=概念 / 2=施策 / 3=詳細" },
    { name = "category",       type = "STRING",  mode = "NULLABLE", description = "concept / entity / claim / evidence / counter" },
    { name = "entity_type",    type = "STRING",  mode = "NULLABLE", description = "organization / person / metric / date (category=entity のみ)" },
    { name = "description",    type = "STRING",  mode = "NULLABLE" },
    { name = "summary_html",   type = "STRING",  mode = "NULLABLE", description = "iframe 向け HTML サマリ。null の場合は description にフォールバック" },
    { name = "source_chunk_id",type = "STRING",  mode = "NULLABLE" },
    { name = "confidence",     type = "FLOAT64", mode = "NULLABLE" },
    { name = "created_at",     type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "DAY"
    field = "created_at"
  }
}

resource "google_bigquery_table" "edges" {
  dataset_id          = google_bigquery_dataset.graph.dataset_id
  table_id            = "edges"
  deletion_protection = true

  schema = jsonencode([
    { name = "document_id",     type = "STRING",  mode = "REQUIRED" },
    { name = "edge_id",         type = "STRING",  mode = "REQUIRED" },
    { name = "source_node_id",  type = "STRING",  mode = "NULLABLE" },
    { name = "target_node_id",  type = "STRING",  mode = "NULLABLE" },
    { name = "edge_type",       type = "STRING",  mode = "NULLABLE", description = "hierarchical / supports / contradicts / related_to / measured_by / involves / causes / exemplifies" },
    { name = "description",     type = "STRING",  mode = "NULLABLE" },
    { name = "weight",          type = "FLOAT64", mode = "NULLABLE" },
    { name = "source_chunk_id", type = "STRING",  mode = "NULLABLE" },
    { name = "created_at",      type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "DAY"
    field = "created_at"
  }
}

resource "google_bigquery_table" "normalization_tools" {
  dataset_id = google_bigquery_dataset.graph.dataset_id
  table_id   = "normalization_tools"

  schema = jsonencode([
    { name = "tool_id",            type = "STRING",  mode = "REQUIRED" },
    { name = "name",               type = "STRING",  mode = "NULLABLE" },
    { name = "version",            type = "STRING",  mode = "NULLABLE" },
    { name = "description",        type = "STRING",  mode = "NULLABLE" },
    { name = "problem_pattern",    type = "STRING",  mode = "NULLABLE", description = "自動マッチング用パターン文字列" },
    { name = "approval_status",    type = "STRING",  mode = "NULLABLE", description = "draft / reviewed / approved / deprecated" },
    { name = "approved_by",        type = "STRING",  mode = "NULLABLE", description = "llm / human" },
    { name = "llm_review_score",   type = "FLOAT64", mode = "NULLABLE", description = "LLM 自動レビュースコア (0〜1)" },
    { name = "llm_review_reason",  type = "STRING",  mode = "NULLABLE" },
    { name = "created_by",         type = "STRING",  mode = "NULLABLE", description = "Firebase Auth UID" },
    { name = "created_at",         type = "TIMESTAMP", mode = "NULLABLE" },
    { name = "updated_at",         type = "TIMESTAMP", mode = "NULLABLE" },
  ])
}

resource "google_bigquery_table" "normalization_tool_runs" {
  dataset_id = google_bigquery_dataset.graph.dataset_id
  table_id   = "normalization_tool_runs"

  schema = jsonencode([
    { name = "run_id",        type = "STRING",  mode = "REQUIRED" },
    { name = "tool_id",       type = "STRING",  mode = "NULLABLE" },
    { name = "document_id",   type = "STRING",  mode = "NULLABLE" },
    { name = "run_type",      type = "STRING",  mode = "NULLABLE", description = "dry_run / apply" },
    { name = "status",        type = "STRING",  mode = "NULLABLE", description = "running / completed / failed" },
    { name = "diff_summary",  type = "STRING",  mode = "NULLABLE", description = "JSON 形式の差分サマリ" },
    { name = "error_message", type = "STRING",  mode = "NULLABLE" },
    { name = "started_at",    type = "TIMESTAMP", mode = "NULLABLE" },
    { name = "completed_at",  type = "TIMESTAMP", mode = "NULLABLE" },
  ])

  time_partitioning {
    type  = "DAY"
    field = "started_at"
  }
}

# plans テーブル (制限値設定。初回 apply 後に bq insert でシードすること)
resource "google_bigquery_table" "plans" {
  dataset_id = google_bigquery_dataset.graph.dataset_id
  table_id   = "plans"

  schema = jsonencode([
    { name = "plan",                      type = "STRING", mode = "REQUIRED", description = "free / pro" },
    { name = "storage_quota_bytes",       type = "INT64",  mode = "NULLABLE" },
    { name = "max_file_size_bytes",       type = "INT64",  mode = "NULLABLE" },
    { name = "max_uploads_per_day",       type = "INT64",  mode = "NULLABLE" },
    { name = "max_members",               type = "INT64",  mode = "NULLABLE" },
    { name = "allowed_extraction_depths", type = "STRING", mode = "NULLABLE", description = "カンマ区切り: full,summary" },
  ])
}

# ---------------------------------------------------------------------------
# Future テーブル (コメントアウト: 後続タスクで有効化)
# ---------------------------------------------------------------------------
# node_aliases, document_topic_mappings, node_scores,
# processing_jobs, graph_snapshots は 10-topic-mapping.md, 11-graph-algorithms.md 参照
