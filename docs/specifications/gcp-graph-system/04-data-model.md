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
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

#### Status Values

- `uploaded`
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
| source_page | INT64 | 元ページ番号 |
| source_offset_start | INT64 | 開始オフセット |
| source_offset_end | INT64 | 終了オフセット |

### nodes

| Column | Type | Description |
| --- | --- | --- |
| document_id | STRING | ドキュメント識別子 |
| node_id | STRING | ノード識別子 |
| label | STRING | 表示ラベル |
| type | STRING | ノード種別 |
| description | STRING | ノード説明 |
| source_chunk_id | STRING | 出典 chunk |
| confidence | FLOAT64 | 生成信頼度 |
| created_at | TIMESTAMP | 作成日時 |

#### Node Type Values

- `abstract`
- `concrete`

### edges

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

- `abstract_to_concrete`
- `related_to`

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

ドキュメントとトピック（abstract ノード）の対応関係を保存する。詳細は [10-topic-mapping.md](10-topic-mapping.md) を参照。

| Column | Type | Description |
| --- | --- | --- |
| mapping_id | STRING | マッピング識別子 |
| document_id | STRING | ドキュメント識別子 |
| topic_node_id | STRING | トピック（abstract ノード）識別子 |
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
