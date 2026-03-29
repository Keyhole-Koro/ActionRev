# 17. BigQuery Migration Strategy

## BigQuery のスキーマ変更制約

BigQuery は RDBMS と異なり、スキーマ変更に厳しい制約がある。

| 操作 | 可否 | 方法 |
|---|---|---|
| カラム追加 (NULLABLE) | ✅ | `ALTER TABLE ADD COLUMN` |
| テーブル追加 | ✅ | `CREATE TABLE IF NOT EXISTS` |
| カラム削除 | ❌ | 削除不可。`_deprecated` サフィックスを付けて放置し、将来テーブル再作成 |
| カラムリネーム | ❌ | 新カラム追加 → バックフィル → 旧カラムを `_deprecated` 化 |
| 型変更 | ❌ (一部△) | 互換性のある変換 (INT64→FLOAT64 など) のみ。基本は新カラム追加 |
| NULLABLE → REQUIRED | ❌ | 不可 |
| REQUIRED → NULLABLE | ✅ | `ALTER TABLE ALTER COLUMN DROP NOT NULL` |

**Terraform の注意点**: `google_bigquery_table` のスキーマを変更すると Terraform は
`deletion_protection = true` で `apply` が失敗する。**初回作成以降のスキーマ変更は
migration scripts で行い、Terraform の schema は変更しない**。

---

## 採用方式: 番号付き SQL ファイル + カスタム Go ランナー

外部ツール (Flyway / Atlas) は導入コストと依存が増えるため、
バックエンドの Go ツールチェーンに統合した軽量カスタムランナーを使用する。

### 管理テーブル

```sql
CREATE TABLE IF NOT EXISTS `{project}.{dataset}.schema_migrations` (
  version     INT64     NOT NULL,
  description STRING,
  applied_at  TIMESTAMP NOT NULL
);
```

### マイグレーションファイル

```
migrations/bigquery/
  000001_baseline.sql          # Terraform で作成済みのテーブルを baseline として記録
  000002_add_node_aliases.sql  # 追加カラム・テーブルはここに書く
  000003_...
```

命名規則: `<6桁連番>_<説明>.sql`

### SQL ファイルの書き方

```sql
-- migrations/bigquery/000002_add_node_aliases.sql
-- description: Add node_aliases table for cross-document deduplication

CREATE TABLE IF NOT EXISTS `{project}.{dataset}.node_aliases` (
  canonical_node_id  STRING NOT NULL,
  alias_node_id      STRING NOT NULL,
  workspace_id       STRING,
  similarity_score   FLOAT64,
  created_at         TIMESTAMP
);
```

プレースホルダー `{project}` と `{dataset}` はランナーが実行時に置換する。

---

## ランナー (`backend/cmd/migrate`)

```
migrate [flags]
  -project  string  GCP プロジェクト ID
  -dataset  string  BigQuery データセット ID (default: graph)
  -dir      string  マイグレーションファイルのディレクトリ (default: migrations/bigquery)
  -dry-run          適用せずに pending を表示するだけ
```

### 実行フロー

```
1. schema_migrations テーブルを CREATE IF NOT EXISTS
2. 適用済みバージョンを SELECT
3. migrations/ の SQL ファイルを昇順で読み込み
4. 未適用のファイルを順に BigQuery で実行
5. 成功ごとに schema_migrations に INSERT
6. 失敗したら即座に停止 (以降の migration は適用しない)
```

### 冪等性

すべての SQL は `CREATE TABLE IF NOT EXISTS` / `ALTER TABLE ADD COLUMN IF NOT EXISTS` を使い
再実行しても安全にする。

---

## CI/CD との統合

デプロイ前に migrate を実行する。

```yaml
# deploy-stage.yml / deploy-prod.yml に追加
- name: Run BigQuery migrations
  run: |
    go run ./backend/cmd/migrate \
      -project=${{ secrets.GCP_PROJECT_ID_STAGE }} \
      -dataset=graph \
      -dir=migrations/bigquery
```

**ロールバック方針**: BigQuery はロールバック DDL を持たない。
問題が起きた場合は前進修正 (forward-only fix) で次の migration として修正を当てる。

---

## カラム削除の手順 (参考)

削除できないため以下の手順を踏む:

1. アプリコードで該当カラムへの読み書きを停止
2. カラム名に `_deprecated_YYYYMMDD` サフィックスを付ける migration を追加 (リネーム不可なので実質コメント扱い)
3. 必要なら `CREATE TABLE ... AS SELECT (削除カラム以外) FROM old_table` でテーブル再作成
4. Terraform の schema も更新 (この場合のみ Terraform schema を変更する)
