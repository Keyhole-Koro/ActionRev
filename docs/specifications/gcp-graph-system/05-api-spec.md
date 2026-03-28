# 05. API Specification

## RPC Contract

### Overview

- フロントエンドとバックエンド間の同期通信は `Connect RPC` を使用する
- 契約定義は `Protocol Buffers` で管理する
- 外部公開 REST API は初期スコープに含めない
- 重い処理は RPC からジョブを起動し、完了確認は status 取得 RPC で行う
- proto の叩き台は [proto/README.md](/home/unix/ActionRev/docs/specifications/gcp-graph-system/proto/README.md) を参照する

## Services

### DocumentService

#### CreateDocument

ファイルをアップロードして document を作成する。

#### Request

- `filename`
- `mime_type`
- `content` またはアップロード対象メタデータ

#### Response

```json
{
  "document_id": "doc_001",
  "status": "uploaded"
}
```

#### Notes

- 実ファイル転送は署名付き URL または別 upload 手段に切り出してもよい

#### GetDocument

document のメタデータと処理状態を取得する。

#### ListDocuments

document 一覧と処理状態を取得する。

#### StartProcessing

document の解析を開始する。

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

可視化用のノード・エッジを取得する。

#### Request Parameters

- `node_type`
- `limit`
- `source_filename` : zip 内の特定ファイル由来のノード・エッジに絞り込む（省略時は全ファイル対象）

#### Response Example

```json
{
  "document_id": "doc_001",
  "nodes": [
    {
      "id": "n1",
      "label": "販売戦略",
      "type": "abstract",
      "description": "販売拡大のための上位方針"
    }
  ],
  "edges": [
    {
      "id": "e1",
      "source": "n1",
      "target": "n2",
      "type": "abstract_to_concrete"
    }
  ]
}
```

### NodeService

#### GetNode

ノード詳細、関連エッジ、出典 chunk を取得する。

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

- service は `DocumentService`, `GraphService`, `NodeService`, `JobService`, `ToolService` に分割する
- request と response は用途単位で明示的に分ける
- `Node`, `Edge`, `Document`, `Job`, `NormalizationTool`, `NormalizationToolRun` は共通 message として定義する
- front でそのまま `React Flow` にマップしやすい field 名を採用する

## Transport Policy

- ブラウザからは `Connect` プロトコルを優先利用する
- 将来的な他クライアント連携に備え、`gRPC` および `gRPC-Web` 互換を維持できる構成を優先する
- 長時間処理は unary RPC で完結させず、job 起動と status 参照に分ける

## Prompt and Extraction Policy

### Prompt Requirements

- 抽象概念と具体概念を明示的に区別させる
- エッジ種別を限定する
- 出典 chunk の参照を必須にする
- JSON Schema に厳密に従うよう要求する

### Post-Processing Requirements

- ラベルの正規化
- 重複ノードの統合
- 不正 JSON の補正または再試行
- 不十分な出力時の fail handling
