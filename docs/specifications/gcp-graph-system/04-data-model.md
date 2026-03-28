# 04. Data Model

## Tables

### documents

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
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

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| node_id | STRING | ノード識別子 |
| label | STRING | 表示ラベル |
| level | INT64 | 階層レベル（0=ドメイン / 1=概念 / 2=施策・アクション / 3=詳細） |
| category | STRING | ノードカテゴリ（`concept` / `entity` / `claim` / `evidence` / `counter`） |
| entity_type | STRING | エンティティ種別（category=entity のみ: `organization` / `person` / `metric` / `date`） |
| description | STRING | ノード説明 |
| summary_html | STRING | ノードサマリの HTML（構造タグのみ、CSS はアプリ側注入）。null の場合は description にフォールバック |
| source_chunk_id | STRING | 出典 chunk |
| confidence | FLOAT64 | 生成信頼度 |
| created_at | TIMESTAMP | 作成日時 |

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

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| edge_id | STRING | エッジ識別子 |
| source_node_id | STRING | 始点ノード |
| target_node_id | STRING | 終点ノード |
| edge_type | STRING | エッジ種別 |
| description | STRING | エッジ説明 |
| weight | FLOAT64 | 重み |
| source_chunk_id | STRING | 出典 chunk |
| created_at | TIMESTAMP | 作成日時 |

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

| Column | Type | Description |
| --- | --- | --- |
| canonical_node_id | STRING | 正規ノード識別子 |
| alias_node_id | STRING | 統合元ノード識別子 |
| alias_label | STRING | 表記揺れラベル |
| similarity_score | FLOAT64 | 類似スコア |
| merge_status | STRING | `suggested` / `approved` / `rejected` |
| created_at | TIMESTAMP | 作成日時 |

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

### processing_jobs

- 非同期ジョブの状態管理に利用する

### graph_snapshots

- 再処理前後の結果比較に利用する

### normalization_tools

- LLM または人手で作成した正規化ツールの定義保存に利用する

### normalization_tool_runs

- ツールの dry-run、本実行、差分、失敗情報の記録に利用する
