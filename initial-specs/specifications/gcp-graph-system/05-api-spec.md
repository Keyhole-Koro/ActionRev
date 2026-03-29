# 05. API Specification

## RPC Contract

### Overview

- フロントエンドとバックエンド間の同期通信は `Connect RPC` を使用する
- 契約定義は `Protocol Buffers` で管理する
- 外部公開 REST API は初期スコープに含めない
- 重い処理は RPC からジョブを起動し、完了確認は status 取得 RPC で行う
- proto の叩き台は [proto/README.md](/home/unix/ActionRev/initial-specs/specifications/gcp-graph-system/proto/README.md) を参照する

## Services

### DocumentService

#### CreateDocument

document を作成し、ファイル upload 用の署名付き URL を発行する。

#### Request

- `filename`
- `mime_type`
- `file_size`

#### Response

```json
{
  "document": {
    "document_id": "doc_001",
    "status": "uploaded"
  },
  "upload_url": "https://storage.googleapis.com/...",
  "upload_method": "PUT",
  "upload_content_type": "application/pdf"
}
```

#### Notes

- 実ファイル転送は `CreateDocument` のレスポンスで返した署名付き URL に対してクライアントが直接実行する
- upload 完了後に `StartProcessing` を呼び出して解析を開始する

#### GetDocument

document のメタデータと処理状態を取得する。

#### ListDocuments

document 一覧と処理状態を取得する。

#### StartProcessing

document の解析を開始する。

#### Preconditions

- 対象 document の実ファイル upload が完了していること
- `documents.status` が `uploaded` であること
- upload 未完了または `processing` / `completed` 状態の document に対してはエラーを返す
- `force_reprocess=true` の場合のみ `completed` または `failed` の document を再処理対象として受け付ける

#### Response

```json
{
  "document_id": "doc_001",
  "status": "processing",
  "job_id": "job_001"
}
```

### GraphService

#### GetGraph

可視化用のノード・エッジを取得する。document 単位の全体表示や初期ロードは `BigQuery` を正本として返す。

#### Request Parameters

- `document_id`
- `workspace_id`
- `category_filters`
- `level_filters`
- `edge_type_filters`
- `limit`
- `source_filename` : zip 内の特定ファイル由来のノード・エッジに絞り込む（省略時は全ファイル対象）
- `resolve_aliases` : canonical ノードへ集約して返すか

#### Response Example

```json
{
  "document_id": "doc_001",
  "nodes": [
    {
      "id": "nd_01JQ8Y7M6Y7YJ8V0X3D4K9P2AB",
      "canonical_node_id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS",
      "scope": "document",
      "label": "販売戦略",
      "level": 1,
      "category": "concept",
      "description": "販売拡大のための上位方針"
    }
  ],
  "edges": [
    {
      "id": "ed_01JQ8YAZE4X7S9N5K2M1P6R8TC",
      "source": "nd_01JQ8Y7M6Y7YJ8V0X3D4K9P2AB",
      "target": "nd_01JQ8Y8T0H4F6B3W9C1N7M2KQD",
      "type": "hierarchical",
      "scope": "document"
    }
  ]
}
```

#### ExpandNeighbors

対話的な近傍探索を行う。canonical node を起点に `Spanner Graph` から指定 hop 数の近傍 subgraph を返す。

#### Request Parameters

- `seed_node_id`
- `workspace_id`
- `max_depth`
- `edge_type_filters`
- `limit_per_hop`
- `resolve_aliases`
- `cross_document`
- `document_ids` : 探索対象を特定 document 群に絞る場合のみ指定

#### Response Example

```json
{
  "graph": {
    "nodes": [
      {"id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "canonical_node_id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "scope": "canonical", "label": "販売戦略"},
      {"id": "cn_01JQ8YD3M5X8F2C7R4T9V1K6LA", "canonical_node_id": "cn_01JQ8YD3M5X8F2C7R4T9V1K6LA", "scope": "canonical", "label": "テレアポ施策"}
    ],
    "edges": [
      {"id": "ed_01JQ8YF1Z6C4N8M2R5T7V9K3LP", "source": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "target": "cn_01JQ8YD3M5X8F2C7R4T9V1K6LA", "type": "hierarchical", "scope": "canonical"}
    ]
  },
  "seed_node_id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS",
  "depth": 1
}
```

#### FindPaths

2 ノード間の多段経路を検索する。`Spanner Graph` を使い、複数経路候補を返す。

#### Request Parameters

- `source_node_id`
- `target_node_id`
- `workspace_id`
- `max_depth`
- `edge_type_filters`
- `limit`
- `cross_document`
- `document_ids`

#### Response Example

```json
{
  "graph": {
    "nodes": [
      {"id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "canonical_node_id": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "scope": "canonical", "label": "販売戦略"},
      {"id": "cn_01JQ8YG6B1N4T8M3R7V2K5P9DX", "canonical_node_id": "cn_01JQ8YG6B1N4T8M3R7V2K5P9DX", "scope": "canonical", "label": "SNS施策"},
      {"id": "cn_01JQ8YJ2F4C6M9T1R3V8K5N7QW", "canonical_node_id": "cn_01JQ8YJ2F4C6M9T1R3V8K5N7QW", "scope": "canonical", "label": "CV率3.2%"}
    ],
    "edges": [
      {"id": "ed_01JQ8YK8M3T6V1R4C7N9P2L5HS", "source": "cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "target": "cn_01JQ8YG6B1N4T8M3R7V2K5P9DX", "type": "hierarchical", "scope": "canonical"},
      {"id": "ed_01JQ8YM4R7C1N5T8V2K6P9L3BZ", "source": "cn_01JQ8YG6B1N4T8M3R7V2K5P9DX", "target": "cn_01JQ8YJ2F4C6M9T1R3V8K5N7QW", "type": "measured_by", "scope": "canonical"}
    ]
  },
  "paths": [
    {
      "node_ids": ["cn_01JQ8YCH9R2V6M4B8T1K5N7PQS", "cn_01JQ8YG6B1N4T8M3R7V2K5P9DX", "cn_01JQ8YJ2F4C6M9T1R3V8K5N7QW"],
      "edge_ids": ["ed_01JQ8YK8M3T6V1R4C7N9P2L5HS", "ed_01JQ8YM4R7C1N5T8V2K6P9L3BZ"],
      "hop_count": 2
    }
  ]
}
```

#### Notes

- `GetGraph` は document 表示用の集約取得を優先する
- `ExpandNeighbors` / `FindPaths` は探索 UX 用であり、低レイテンシを優先して `Spanner Graph` を参照する
- `BigQuery` と `Spanner Graph` に同期遅延がある場合、探索結果は最新の抽出完了直後とわずかにずれる可能性がある
- 探索系 RPC は必ず `workspace_id` を受け取り、workspace 境界をまたがる探索は許可しない
- `cross_document=false` の場合は現在の document または `document_ids` の範囲だけを探索対象とする
- `cross_document=true` の場合は同一 workspace 内の canonical graph を探索対象とする
- `Node.scope=document` の場合 `id` は `nd_*` を返し、`canonical_node_id` は alias 解決済みなら補助属性として返す
- `Node.scope=canonical` の場合 `id` と `canonical_node_id` は同一の `cn_*` を返し、`document_id` は必須ではない
- `Edge.scope=document` の場合 `source` / `target` は document node (`nd_*`) を指す
- `Edge.scope=canonical` の場合 `source` / `target` は canonical node (`cn_*`) を指す

### NodeService

#### GetNode

ノード詳細、関連エッジ、出典 chunk を取得する。

#### Request Parameters

- `document_id`
- `node_id`
- `resolve_aliases` : alias ノード指定時に canonical ノードへ寄せて返すか

### JobService

#### GetJobStatus

処理ジョブの状態を取得する。

#### Response Example

```json
{
  "job_id": "job_001",
  "status": "running"
}
```

### ToolService

`dev` ロールを持つ管理者のみアクセス可能。workspace とは無関係のシステムグローバルなツール管理 API。詳細は [09-normalization-tools.md](09-normalization-tools.md) を参照。

#### GenerateNormalizationTool

問題パターンの説明やサンプルデータをもとに、LLM から Python 正規化スクリプト案を生成する。

#### SaveNormalizationTool

生成されたスクリプトを `problem_pattern` と manifest とともに保存する（`draft` 状態）。

#### ListNormalizationTools

正規化ツール一覧を取得する。`approval_status` でフィルタ可能。

#### UpdateNormalizationToolStatus

ツールの状態を遷移させる（`draft` → `reviewed` → `approved` / `deprecated`）。`approved` 状態のみ本番適用可能。

#### RunNormalizationTool

ツールをサンドボックスで dry-run または本実行する。`APPLY` モードは `approved` のツールのみ実行可能。

#### GetNormalizationToolRun

ツール実行結果、差分、ログ、出力物参照を取得する。

## Proto Design Guidelines

- package は単一の `actionrev.graph.v1` とし、versioning は package suffix で管理する
- `.proto` ファイルは service 単位で分割し、1ファイル1service を原則とする
- service は `UserService`, `WorkspaceService`, `DocumentService`, `GraphService`, `NodeService`, `JobService`, `ToolService` に分割する
- `GraphService` はグラフ取得と開発者向け統計 RPC を持つ
- `GraphService` は `GetGraph` のほか、`ExpandNeighbors` と `FindPaths` の探索 RPC を持つ
- request と response は用途単位で明示的に分ける
- `Node`, `Edge`, `Document`, `Job`, `NormalizationTool`, `NormalizationToolRun` などの共通 message / enum は `common.proto` に集約する
- 複数 service から参照される message は service 個別 proto に重複定義しない
- package を domain ごとに細分化するのは初期スコープ外とし、import と生成コードの複雑化を避ける
- front でそのまま `React Flow` にマップしやすい field 名を採用する
- ノード分類は `level` / `category` / `entity_type` を正とし、旧2値分類は持ち込まない

### Proto File Ownership

- `common.proto`: 共通 message / enum のみを保持し、service は定義しない
- `document.proto`: upload 開始と document メタデータ取得のみを扱う
- `graph.proto`: 可視化用グラフ取得と `/dev/stats` 系 RPC を扱う
- `node.proto`: 単一ノード詳細取得を扱う
- `job.proto`: 非同期ジョブ状態取得のみを扱う
- `tool.proto`: 正規化ツールの生成、承認、実行を扱う
- `user.proto`, `workspace.proto`: 認証後のユーザー同期と workspace 管理を扱う

### Package Evolution Policy

- 後方互換を壊す変更は `actionrev.graph.v2` を新設して行う
- `v1` では field 追加を許容し、field 削除・型変更・意味変更は禁止する
- `buf breaking` を導入した時点で `main` ブランチとの差分を自動検証する

## Transport Policy

- ブラウザからは `Connect` プロトコルを優先利用する
- 将来的な他クライアント連携に備え、`gRPC` および `gRPC-Web` 互換を維持できる構成を優先する
- 長時間処理は unary RPC で完結させず、job 起動と status 参照に分ける
- 探索系 RPC は 1 画面あたり少数回の集約呼び出しを前提とし、N+1 的な node 単位 fetch を避ける

## Prompt and Extraction Policy

### Prompt Requirements

- ノードの `level` と `category` を明示的に割り当てさせる
- `level` は常に 0〜3 の固定4段階で割り当て、文書ごとに段数を変えさせない
- エッジ種別を限定する
- 出典 chunk の参照を必須にする
- JSON Schema に厳密に従うよう要求する

### Post-Processing Requirements

- ラベルの正規化
- 重複ノードの統合
- 不正 JSON に対する JSON repair を 1 回だけ試行する
- JSON repair 後も不正な場合は Gemini 再試行を最大 2 回まで行う
- JSON repair は syntax error のみを対象とし、semantic error は補正しない
- semantic error は schema 必須項目欠落、enum 不正値、参照不整合、制約違反を含む
- 構造成立に必須な項目はフォールバックせず、要素破棄または再試行対象とする
- `description`, `summary_html`, `entity_type` など品質補助項目に限ってフォールバックを許容する
- 不十分な出力時の fail handling
- chunk 抽出の確定失敗は document 全体の失敗として扱う
