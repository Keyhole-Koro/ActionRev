# 11. Graph Algorithms

## Overview

トピックマッピングと横断グラフ構築が完了した後、グラフアルゴリズムを適用して概念の重要度・クラスタ・関係経路を分析する。アルゴリズム処理は Cloud Run Jobs でバッチ実行し、結果を BigQuery に保存してフロントエンドの可視化に利用する。

## 対象アルゴリズム

| アルゴリズム | 用途 | 優先度 |
| --- | --- | --- |
| PageRank / 中心性分析 | 重要概念の自動ランキング | 高 |
| コミュニティ検出 | トピッククラスタの可視化 | 高 |
| 最短経路 | 概念間の関係経路の探索 | 中 |
| 類似ノード推薦 | 関連概念のサジェスト | 中 |
| 到達可能性分析 | 影響範囲の把握 | 低 |

## 処理アーキテクチャ

```
[BigQuery: nodes / edges]
         ↓ 処理完了トリガー or 定期実行
[Cloud Run Jobs: graph algo worker]
  - NetworkX または igraph で計算
  - PageRank スコア / コミュニティ ID / 中心性を計算
         ↓
[BigQuery: node_scores テーブル]
  - node_id, algo_type, score, computed_at
```

アルゴリズム処理をメインの API サーバから分離することで、重い計算がリクエスト処理に影響しない。

## データモデル

### node_scores

| Column | Type | Description |
| --- | --- | --- |
| node_id | STRING | ノード識別子 |
| algo_type | STRING | アルゴリズム種別 |
| score | FLOAT64 | スコア値 |
| metadata | JSON | アルゴリズム固有の追加情報 |
| computed_at | TIMESTAMP | 計算日時 |

#### algo_type Values

- `pagerank` : PageRank スコア
- `degree_centrality` : 次数中心性
- `betweenness_centrality` : 媒介中心性
- `community_id` : コミュニティ検出結果（metadata に community ラベルを含む）

### node_aliases（オントロジー統合との共有）

| Column | Type | Description |
| --- | --- | --- |
| canonical_node_id | STRING | 正規ノード識別子 |
| alias_node_id | STRING | 統合元ノード識別子 |
| alias_label | STRING | 表記揺れラベル |
| similarity_score | FLOAT64 | 類似スコア |
| merge_status | STRING | 統合状態 |
| created_at | TIMESTAMP | 作成日時 |

#### merge_status Values

- `suggested` : 自動候補（未レビュー）
- `approved` : 承認済み（グラフクエリ時に canonical に読み替える）
- `rejected` : 却下

## グラフエンジンの選択

### 初期：BigQuery + Cloud Run Jobs

- BigQuery の再帰 CTE で最短経路・到達可能性を実装する
- NetworkX / igraph を Cloud Run Jobs で動かし、計算結果を BQ に書き戻す
- 実装コストが低く、既存スタックから外れない

```sql
-- BigQuery での多段経路の例（再帰 CTE）
WITH RECURSIVE paths AS (
  SELECT source_node_id, target_node_id, 1 AS depth
  FROM edges
  WHERE source_node_id = @start_node_id
  UNION ALL
  SELECT e.source_node_id, e.target_node_id, p.depth + 1
  FROM edges e
  JOIN paths p ON e.source_node_id = p.target_node_id
  WHERE p.depth < 5
)
SELECT * FROM paths
```

### 将来：Spanner Graph

- GQL（Graph Query Language）ネイティブ対応で多段エッジトラバーサルが直感的に書ける
- `MATCH (n)-[e*1..5]->(m)` のような経路指定が可能
- BigQuery との二重持ちになるが、クエリ特性が大きく異なるため大規模グラフでは有効
- BQ をコールドデータ・分析用途、Spanner Graph をリアルタイム探索用途で使い分ける

```
-- Spanner Graph GQL の例
GRAPH ActionRevGraph
MATCH (start {id: @start_node_id})-[*1..5]->(m)
RETURN m.label, m.type
```

## フロントエンドへの反映

- `GetGraph` のレスポンスに `node_scores` を JOIN して返す
- スコアに応じてノードのサイズ・色・強調表示を変える
- コミュニティ ID でノードをグループ色分けする

## 実行タイミング

| タイミング | 用途 |
| --- | --- |
| document 処理完了後に即時実行 | document 単位の中心性・PageRank 更新 |
| 夜間バッチ（BigQuery scheduled query）| 横断グラフ全体の再計算 |
| 手動トリガー | 再処理・デバッグ |

## 処理フローの全体像

```
[正規化 → 抽出]
      ↓
[ノード統合（document 内）]
      ↓
[トピックマッピング（heuristic → LLM）]
      ↓
[トピック canonical 化（node_aliases）]
      ↓
[横断グラフ構築]
      ↓
[グラフアルゴリズム適用（Cloud Run Jobs）]
      ↓
[可視化（React Flow）]
```

## Open Issues

- Cloud Run Jobs のトリガー方法（document 処理完了後の Pub/Sub か、定期実行か）
- Spanner Graph への移行タイミングと BQ との使い分け方針
- アルゴリズムの再計算スコープ（全グラフか差分か）
- フロントでのスコア可視化の UX（サイズ変化・色分け・ツールチップ）
