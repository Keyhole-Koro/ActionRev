# 10. Topic Mapping

## Overview

複数ドキュメントにまたがる概念統合を実現するため、抽象ノードをトピックとして扱い、各ドキュメントのノード群をトピックにマッピングする。マッピングはヒューリスティックによる候補絞り込みと LLM による検証の二段階で行い、LLM 呼び出しコストを抑えながら精度を確保する。

## トピックの定義

- 既存データモデルの `abstract` ノードがトピックに相当する
- トピックは単一ドキュメント内で生成されるが、複数ドキュメントから参照される共有概念として扱う
- トピック間の重複は `node_aliases` テーブルでカノニカル化する（詳細は 11-graph-algorithms.md を参照）

## マッピングの処理フロー

```
Step 1: ヒューリスティックで候補を絞り込む
  - chunk/node の embedding とトピック embedding の cosine similarity を計算する
  - 閾値（初期値 0.80）を超えたペアを候補として残す
  - LLM に投げる件数を大幅に削減する

Step 2: LLM で候補を検証する
  - 候補ペアについて「このノード群はこのトピックに属するか」を判定させる
  - 判定理由を必ず返させ、監査可能にする

Step 3: マッピング結果を保存する
  - document_topic_mappings テーブルに保存する
  - method フィールドで heuristic / llm を記録する
```

## ヒューリスティック手法

| 手法 | 精度 | コスト | 用途 |
| --- | --- | --- | --- |
| キーワードマッチ | 低〜中 | ほぼ無料 | 粗い一次絞り込み |
| 編集距離（ラベル類似度） | 中 | ほぼ無料 | 表記揺れの検出 |
| Embedding cosine similarity | 高 | Vertex AI Embeddings 呼び出し | 二次ランキング |

初期実装はキーワードマッチと編集距離で粗く絞り込み、embedding で再ランキングし、最後に LLM で確定する三段階構成を推奨する。

## データモデル

### document_topic_mappings

| Column | Type | Description |
| --- | --- | --- |
| mapping_id | STRING | マッピング識別子 |
| document_id | STRING | ドキュメント識別子 |
| topic_node_id | STRING | トピック（abstract ノード）識別子 |
| confidence | FLOAT64 | 信頼スコア |
| reason | STRING | LLM による判定理由 |
| method | STRING | マッピング手法 |
| created_at | TIMESTAMP | 作成日時 |

#### Method Values

- `keyword` : キーワードマッチによる自動マッピング
- `embedding` : embedding 類似度による自動マッピング
- `llm` : LLM による検証済みマッピング
- `manual` : 人手によるマッピング

## トピックのカノニカル化

ドキュメントをまたいで同一概念が異なるラベルで生成される場合がある（例：「販売戦略」と「セールス戦略」）。このため `node_aliases` テーブルと組み合わせてトピックを統合する。

```
[各 doc の abstract ノード生成]
         ↓
[embedding similarity で類似トピックを候補化]
         ↓
[LLM で同一概念かを判定]
         ↓
[node_aliases に canonical_node_id として登録]
         ↓
[GetGraph 時に canonical ノードに集約して返す]
```

GraphService の `GetGraph` に `resolve_aliases=true` オプションを追加し、BQ クエリ側で alias JOIN を行う。

## 横断グラフの構造

```
トピック A（canonical abstract ノード）
  ├─ doc_001 の concrete ノード群
  ├─ doc_002 の concrete ノード群
  └─ doc_003 の concrete ノード群

トピック B（canonical abstract ノード）
  ├─ doc_001 の concrete ノード群
  └─ doc_004 の concrete ノード群
```

## API 拡張

### GraphService への追加

- `GetTopicMap` : トピック一覧と各トピックに紐づくドキュメント・ノード数を返す
- `GetGraph` の `resolve_aliases` パラメータ追加

### TopicService（将来）

- `ListTopics` : トピック一覧と統計情報を返す
- `GetTopicDocuments` : トピックに紐づくドキュメント一覧を返す
- `MergeTopics` : 人手によるトピック統合を実行する

## Open Issues

- embedding の生成タイミング（抽出直後か、バッチか）
- cosine similarity の閾値の初期値と調整方針
- LLM に渡すコンテキスト量（ノード label だけか description も含めるか）
- カノニカル化の承認フローをどこまで厳密にするか
