# 04. Data Model

## Tables

### users

Firebase Auth でログインした際に初回のみ自動作成する。

| Column | Type | Description |
| --- | --- | --- |
| user_id | STRING | Firebase Auth UID |
| email | STRING | メールアドレス |
| display_name | STRING | 表示名 |
| created_at | TIMESTAMP | 初回ログイン日時 |
| last_login_at | TIMESTAMP | 最終ログイン日時 |

### user と workspace の関係

```
User
  │ 1
  │ ├─ 複数の workspace を作成できる（owner）
  │ │      workspaces.owner_id = user_id
  │ │
  │ └─ 複数の workspace にメンバーとして参加できる
  │        workspace_members.user_id = user_id
  │ *
Workspace
  │ 1
  └─ 複数の document を持つ
         documents.workspace_id = workspace_id
```

- 1ユーザーは複数の workspace を作成できる
- 1ユーザーは複数の workspace にメンバーとして参加できる
- 1つの workspace は複数の document を持つ
- workspace の作成者は自動で `editor` として `workspace_members` に登録される

### workspaces

| Column | Type | Description |
| --- | --- | --- |
| workspace_id | STRING | ワークスペース識別子 |
| name | STRING | ワークスペース名 |
| owner_id | STRING | オーナーのユーザーID（Firebase Auth UID） |
| plan | STRING | `free` / `pro` |
| stripe_customer_id | STRING | Stripe の顧客ID |
| stripe_subscription_id | STRING | Stripe のサブスクリプションID |
| storage_used_bytes | INT64 | 使用済みストレージ容量 |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### workspace_members

| Column | Type | Description |
| --- | --- | --- |
| workspace_id | STRING | ワークスペース識別子 |
| user_id | STRING | メンバーのユーザーID |
| role | STRING | `editor` / `viewer` / `dev` |
| invited_at | TIMESTAMP | 招待日時 |

#### Role Values

- `editor` : ドキュメントのアップロード・削除・処理実行が可能
- `viewer` : グラフの閲覧のみ可能
- `dev` : `/dev/stats` アクセス可能（editor/viewer に追加付与）

### documents

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| workspace_id | STRING | 所属ワークスペース識別子 |
| uploaded_by | STRING | アップロードしたユーザーID |
| filename | STRING | 元ファイル名 |
| gcs_uri | STRING | 保存先 URI |
| mime_type | STRING | MIME type |
| file_size | INT64 | ファイルサイズ |
| status | STRING | 処理状態 |
| extraction_depth | STRING | 抽出粒度（`full` / `summary`）。デフォルトは `full` |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

#### Status Values

- `uploaded` : メタデータ登録とファイル upload 完了後、解析開始前
- `pending_normalization` : 正規化ツールの承認待ちで処理を停止中
- `processing`
- `completed`
- `failed`

### document_chunks

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| chunk_id | STRING | chunk 識別子 |
| chunk_index | INT64 | chunk 順序 |
| text | STRING | chunk テキスト |
| source_filename | STRING | 元ファイル名（zip 内ファイルの場合は展開後のファイル名、単ファイルの場合は document の filename と同値） |
| source_page | INT64 | 元ページ番号 |
| source_offset_start | INT64 | 開始オフセット |
| source_offset_end | INT64 | 終了オフセット |

### nodes

詳細な抽出戦略は [12-extraction-strategy.md](12-extraction-strategy.md) を参照。

`node_id` は表示ラベルや LLM が返す `local_id` から直接導出せず、永続化時にバックエンドがグローバル一意な値を採番する。初期実装では `nd_<ULID>` 形式を推奨する。

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| node_id | STRING | ノード識別子 |
| extraction_local_id | STRING | LLM 抽出時のローカル識別子（document 内で追跡用） |
| label | STRING | 表示ラベル |
| level | INT64 | 階層レベル（0=ドメイン / 1=概念 / 2=施策・アクション / 3=詳細） |
| category | STRING | ノードカテゴリ（`concept` / `entity` / `claim` / `evidence` / `counter`） |
| entity_type | STRING | エンティティ種別（category=entity のみ: `organization` / `person` / `metric` / `date` / `location`） |
| description | STRING | ノード説明 |
| summary_html | STRING | ノードサマリの HTML（構造タグのみ、CSS はアプリ側注入）。null の場合は description にフォールバック |
| source_chunk_id | STRING | 出典 chunk |
| confidence | FLOAT64 | 生成信頼度 |
| created_at | TIMESTAMP | 作成日時 |

- API 返却時の `canonical_node_id` は `node_aliases` を解決して補完する派生属性であり、`nodes` テーブルの永続カラムには含めない
- `GetGraph` の node は `id = node_id`, `scope = document` を返す
- `ExpandNeighbors` / `FindPaths` の node は `id = canonical_node_id`, `scope = canonical` を返す

#### Node Category Values

- `concept` : 抽象的・具体的な概念
- `entity` : 実体（組織・人物・数値・日付）
- `claim` : 主張・判断・結論
- `evidence` : 主張を支持する根拠・事例
- `counter` : 主張への反論・留意点

#### Node Level Values

- `0` : ドメイン（文書全体を覆う最上位概念）
- `1` : 概念（主要なテーマ・方針）
- `2` : 施策・アクション（具体的な取り組み）
- `3` : 詳細（数値・固有名詞・具体的事実）

### edges

詳細な抽出戦略は [12-extraction-strategy.md](12-extraction-strategy.md) を参照。

`edge_id` も `source_node_id` / `target_node_id` / `edge_type` の連結値ではなく、永続化時にバックエンドがグローバル一意な値を採番する。初期実装では `ed_<ULID>` 形式を推奨する。

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| edge_id | STRING | エッジ識別子 |
| extraction_local_id | STRING | LLM 抽出時のローカル識別子（document 内で追跡用、任意） |
| source_node_id | STRING | 始点ノード |
| target_node_id | STRING | 終点ノード |
| edge_type | STRING | エッジ種別 |
| description | STRING | エッジ説明 |
| weight | FLOAT64 | 重み |
| source_chunk_id | STRING | 出典 chunk |
| created_at | TIMESTAMP | 作成日時 |

- API 返却時の `Edge.scope` は派生属性であり、`edges` テーブルの永続カラムには含めない
- `GetGraph` の edge は `scope = document` を返す
- `ExpandNeighbors` / `FindPaths` の edge は `scope = canonical` を返す

#### Edge Type Values

- `hierarchical` : 階層の親子関係
- `supports` : evidence が claim を支持する
- `contradicts` : counter が claim に反論する
- `related_to` : 汎用的な関連
- `measured_by` : concept/entity が metric で測定される
- `involves` : concept に entity が関与する
- `causes` : concept が別の concept を引き起こす
- `exemplifies` : 上位概念の具体例

## Future Tables

### node_aliases

トピックのカノニカル化とオントロジー統合に利用する。詳細は [10-topic-mapping.md](10-topic-mapping.md) を参照。

`canonical_node_id` は文書内 `node_id` とは別の識別空間として扱う。初期実装では `cn_<ULID>` 形式を推奨する。

| Column | Type | Description |
| --- | --- | --- |
| workspace_id | STRING | 所属ワークスペース識別子 |
| canonical_node_id | STRING | 正規ノード識別子 |
| alias_node_id | STRING | 統合元ノード識別子 |
| alias_label | STRING | 表記揺れラベル |
| similarity_score | FLOAT64 | 類似スコア |
| merge_status | STRING | `suggested` / `approved` / `rejected` |
| created_at | TIMESTAMP | 作成日時 |

### canonical_nodes

探索用 graph と document 横断 UI の主キーとして使う canonical node の属性を保持する。`Spanner Graph` への同期元になる。

| Column | Type | Description |
| --- | --- | --- |
| workspace_id | STRING | 所属ワークスペース識別子 |
| canonical_node_id | STRING | 正規ノード識別子 |
| label | STRING | canonical 表示ラベル |
| category | STRING | 代表カテゴリ |
| level_hint | INT64 | 代表 level |
| description | STRING | canonical 説明 |
| representative_node_id | STRING | 代表元 node_id |
| updated_at | TIMESTAMP | 更新日時 |

- `canonical_nodes.canonical_node_id` は探索 API の `Node.id` として露出する
- document ノードと canonical ノードの対応は `node_aliases.alias_node_id -> canonical_node_id` で追跡する

### document_topic_mappings

ドキュメントとトピック（`category=concept` かつ `level in (0, 1)` の canonical ノード）の対応関係を保存する。詳細は [10-topic-mapping.md](10-topic-mapping.md) を参照。

| Column | Type | Description |
| --- | --- | --- |
| mapping_id | STRING | マッピング識別子 |
| document_id | STRING | ドキュメント識別子 |
| topic_node_id | STRING | トピックノード識別子 |
| confidence | FLOAT64 | 信頼スコア |
| reason | STRING | LLM による判定理由 |
| method | STRING | `keyword` / `embedding` / `llm` / `manual` |
| created_at | TIMESTAMP | 作成日時 |

### node_scores

グラフアルゴリズムの計算結果を保存する。詳細は [11-graph-algorithms.md](11-graph-algorithms.md) を参照。

| Column | Type | Description |
| --- | --- | --- |
| node_id | STRING | ノード識別子 |
| algo_type | STRING | `pagerank` / `degree_centrality` / `betweenness_centrality` / `community_id` |
| score | FLOAT64 | スコア値 |
| metadata | JSON | アルゴリズム固有の追加情報 |
| computed_at | TIMESTAMP | 計算日時 |

### graph_sync_jobs

`BigQuery` 正本から `Spanner Graph` への同期状態を管理する。

| Column | Type | Description |
| --- | --- | --- |
| sync_job_id | STRING | 同期ジョブ識別子 |
| workspace_id | STRING | 対象ワークスペース |
| document_id | STRING | 対象 document（全体同期時は null 可） |
| status | STRING | `queued` / `running` / `completed` / `failed` |
| synced_node_count | INT64 | 同期ノード数 |
| synced_edge_count | INT64 | 同期エッジ数 |
| started_at | TIMESTAMP | 開始日時 |
| completed_at | TIMESTAMP | 完了日時 |

### processing_jobs

- 非同期ジョブの状態管理に利用する

### graph_snapshots

- 再処理前後の結果比較に利用する

### plans

プランごとの制限値を管理する設定テーブル。

| Column | Type | Description |
| --- | --- | --- |
| plan | STRING | `free` / `pro` |
| storage_quota_bytes | INT64 | ストレージ上限 |
| max_file_size_bytes | INT64 | 1ファイルあたりの上限サイズ |
| max_uploads_per_day | INT64 | 1日あたりのアップロード上限 |
| max_members | INT64 | workspace メンバー上限 |
| allowed_extraction_depths | STRING | 使用可能な extraction_depth（カンマ区切り） |

#### デフォルト値

| | free | pro |
| --- | --- | --- |
| storage_quota_bytes | 1GB | 50GB |
| max_file_size_bytes | 50MB | 500MB |
| max_uploads_per_day | 10 | 200 |
| max_members | 3 | 20 |
| allowed_extraction_depths | `summary` | `full,summary` |

### normalization_tools

| Column | Type | Description |
| --- | --- | --- |
| tool_id | STRING | ツール識別子 |
| name | STRING | ツール名 |
| version | STRING | バージョン |
| description | STRING | 説明 |
| problem_pattern | STRING | 対処する問題パターン（自動マッチング用） |
| approval_status | STRING | `draft` / `reviewed` / `approved` / `deprecated` |
| approved_by | STRING | `llm` / `human` |
| llm_review_score | FLOAT64 | LLM 自動レビューの信頼スコア（0〜1） |
| llm_review_reason | STRING | LLM の判定理由 |
| created_by | STRING | 作成者ユーザーID |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

### normalization_tool_runs

- ツールの dry-run、本実行、差分、失敗情報の記録に利用する
