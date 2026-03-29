# 12. Extraction Strategy

## Overview

高精度な知識グラフを構築するため、以下の4つの手法を組み合わせて抽出を行う。

- **セマンティックチャンキング**: LLM が意味の切れ目でドキュメントを分割する
- **多階層ノード**: レベルとカテゴリで表現する
- **2パス抽出**: chunk 単位の細粒度抽出と、文書全体の統合・階層化を分離する
- **クレーム／エビデンス型 + エンティティ型**: 概念だけでなく論理構造と実体を抽出する

---

## Stage 3: Semantic Chunking

固定サイズ分割の代わりに、LLM がドキュメントの意味的な区切りを判断してチャンクを生成する。

### 入力

- 正規化済みドキュメントの全文（または分割可能なセクション単位）

### Gemini への指示方針

- セクション・段落・論点の切れ目を認識させる
- 1チャンクは「1つのトピックまたは論点を扱う単位」とする
- チャンクサイズの上限（例: 2000トークン）を設け、超える場合はさらに分割する

### 出力

```json
{
  "chunks": [
    {
      "chunk_index": 0,
      "heading": "背景と課題",
      "text": "..."
    },
    {
      "chunk_index": 1,
      "heading": "施策A: テレアポ強化",
      "text": "..."
    }
  ]
}
```

---

## Stage 4: Pass 1 — Fine-grained Extraction（chunk 単位）

各チャンクに対して Gemini を呼び出し、細粒度で全要素を抽出する。

### 抽出対象

| category | 説明 | 例 |
| --- | --- | --- |
| `concept` | 抽象的・具体的な概念 | 販売戦略、テレアポ施策 |
| `entity` | 実体（組織・人物・数値・日付） | A社、CV率3.2%、2026年Q1 |
| `claim` | 主張・判断・結論 | "SNSの方がROIが高い" |
| `evidence` | 主張を支持する根拠・事例 | "A社でCV率3.2%を達成" |
| `counter` | 主張への反論・留意点 | "テレアポは関係構築に強み" |

### エッジ抽出対象

| edge_type | 説明 |
| --- | --- |
| `supports` | evidence → claim |
| `contradicts` | counter → claim |
| `related_to` | 汎用的な関連 |
| `measured_by` | concept/entity → metric entity |
| `involves` | concept → entity |
| `causes` | concept → concept |
| `exemplifies` | 上位概念 → 具体例 |

### 出力スキーマ

```json
{
  "nodes": [
    {
      "local_id": "n1",
      "label": "テレアポ施策",
      "category": "concept",
      "level": 2,
      "entity_type": null,
      "description": "...",
      "source_chunk_id": "c_001"
    },
    {
      "local_id": "n2",
      "label": "CV率 3.2%",
      "category": "entity",
      "level": 3,
      "entity_type": "metric",
      "description": "...",
      "source_chunk_id": "c_001"
    }
  ],
  "edges": [
    {
      "source": "n1",
      "target": "n2",
      "edge_type": "measured_by",
      "description": "...",
      "source_chunk_id": "c_001"
    }
  ]
}
```

### フィールド要件（Pass 1）

| フィールド | 必須 | 欠落時の扱い |
| --- | --- | --- |
| `nodes` | 必須 | 配列自体がなければ chunk 全体を再試行対象とする |
| `edges` | 必須 | 配列自体がなければ空配列には補正せず、chunk 全体を再試行対象とする |
| `nodes[].local_id` | 必須 | 欠落した node は破棄する |
| `nodes[].label` | 必須 | 欠落した node は破棄する |
| `nodes[].category` | 必須 | 不正値を含む node は破棄する |
| `nodes[].level` | 必須 | 0〜3 以外の値を含む node は破棄する |
| `nodes[].source_chunk_id` | 必須 | 欠落した node は破棄する |
| `nodes[].entity_type` | 条件付き必須 | `category=entity` で欠落または不正値なら `unspecified` として保存する |
| `nodes[].description` | 任意 | 欠落時は空として保存する |
| `edges[].source` | 必須 | 欠落した edge は破棄する |
| `edges[].target` | 必須 | 欠落した edge は破棄する |
| `edges[].edge_type` | 必須 | 不正値を含む edge は破棄する |
| `edges[].source_chunk_id` | 必須 | 欠落した edge は破棄する |
| `edges[].description` | 任意 | 欠落時は空として保存する |

- `nodes[]` / `edges[]` の配列要素単位で破棄可能な不正値は repair せず、その要素のみ破棄する
- 破棄後に `nodes[]` が空になった場合は chunk 全体を semantic failure とみなし、再試行対象とする
- `edges[].source` / `target` が存在しない `local_id` を参照している edge は破棄する

---

## Stage 5: Pass 2 — Document-level Synthesis（文書全体）

Pass 1 の全チャンク抽出結果をまとめて Gemini に投入し、文書全体の構造を把握させる。

### 処理内容

1. **重複統合**: 同一概念の表記揺れを統合し、canonical ラベルを決定する
2. **階層割り当て**: 各ノードに level（0〜3）を付与する
3. **クレーム構造の整理**: claim / evidence / counter の論理関係を明確化する
4. **上位概念の補完**: level 0〜1 の抽象概念が不足している場合は補完する
5. **エッジの補完・整理**: chunk をまたぐ関係（`hierarchical`, `causes` など）を追加する

### 階層レベルの定義

| level | 名称 | 説明 | 例 |
| --- | --- | --- | --- |
| 0 | ドメイン | 文書全体を覆う最上位概念 | 事業戦略 |
| 1 | 概念 | 主要なテーマ・方針 | 販売戦略、マーケティング戦略 |
| 2 | 施策・アクション | 具体的な取り組み | テレアポ施策、SNS施策 |
| 3 | 詳細 | 数値・固有名詞・具体的事実 | CV率3.2%、スクリプト改善 |

### level 割り当てルール

- `level` は文書ごとに段数を変えず、常に 0〜3 の4段階で割り当てる
- `level=0` は文書全体を一言で束ねる最上位テーマに限定し、通常は 0 件または 1 件、多くても少数に抑える
- `level=1` は文書の主要なテーマ、方針、章レベルの概念に割り当てる
- `level=2` は具体的な施策、アクション、ワークストリーム、実行手段に割り当てる
- `level=3` は固有名詞、数値、日付、具体的事実、事例、補足ディテールに割り当てる
- `category=entity` は原則として `level=3` に割り当てる
- `category=claim` / `evidence` / `counter` は内容に応じて `level=2` または `level=3` を選び、文書全体テーマに相当しない限り `level=0` にしない
- 親子関係が明確な場合、子ノードは親ノードより下位の `level` に割り当てる
- 同一文書内では、同等の抽象度を持つノードに同じ `level` を割り当て、相対評価ではなく役割ベースで判断する
- `hierarchical` は近い `level` 間を優先し、不要な飛び級を避ける

### 出力スキーマ

Pass 1 と同形式で、`node_id`（文書内で確定した ID）を付与して返す。

### フィールド要件（Pass 2）

| フィールド | 必須 | 欠落時の扱い |
| --- | --- | --- |
| `nodes` | 必須 | 配列自体がなければ document 全体を再試行対象とする |
| `edges` | 必須 | 配列自体がなければ空配列には補正せず、document 全体を再試行対象とする |
| `nodes[].node_id` | 必須 | 欠落した node は破棄する |
| `nodes[].label` | 必須 | 欠落した node は破棄する |
| `nodes[].category` | 必須 | 不正値を含む node は破棄する |
| `nodes[].level` | 必須 | 0〜3 以外の値を含む node は破棄する |
| `nodes[].description` | 任意 | 欠落時は Pass 1 の description を引き継ぐ |
| `nodes[].entity_type` | 条件付き必須 | `category=entity` で欠落または不正値なら `unspecified` として保存する |
| `edges[].source` | 必須 | 欠落した edge は破棄する |
| `edges[].target` | 必須 | 欠落した edge は破棄する |
| `edges[].edge_type` | 必須 | 不正値を含む edge は破棄する |
| `edges[].description` | 任意 | 欠落時は空として保存する |

- Pass 2 では document 構造の成立を優先し、要素単位の不正はその node / edge のみ破棄する
- 破棄後に level 0〜2 の構造ノードが全て消える場合は semantic failure とみなし、document 全体を再試行対象とする
- `description` が欠落しても canonical 化や可視化の主処理は継続し、HTML summary 生成で補完機会を持つ

---

## Extraction Depth

抽出の粒度は `extraction_depth` パラメータで切り替え可能とする。`StartProcessing` RPC で指定する。

| 値 | 説明 | 対象 level |
| --- | --- | --- |
| `full` | 数値・固有名詞レベルまで全て抽出する | 0〜3 |
| `summary` | 施策・アクションまで抽出する（詳細は親ノードの description に含める） | 0〜2 |

- デフォルトは `full`
- Pass 2 で `extraction_depth=summary` の場合、level 3 ノードを親ノードに統合して削除する
- `documents` テーブルに使用した `extraction_depth` を記録し、再処理時に参照可能にする

---

## Context Injection Policy

入力トークンは安価であるため、各ステージで出力精度を最大化するためにコンテキストを積極的に注入する。

### Layer 1: 全ステージ共通（常時注入）

- セマンティックチャンキングで生成した文書アウトライン（heading 一覧）

### Layer 2: ステージ別注入

| ステージ | 追加注入するコンテキスト |
| --- | --- |
| Pass 1（chunk N 処理時） | 全チャンクテキスト + 処理対象 chunk N の明示 |
| HTML サマリ生成 | 対象ノードの隣接ノード（親・子・関連）+ 出典チャンクの原文 |

### Layer 3: 横断注入（全ステージ）

- 他ドキュメントの level 0〜1 ノード（topic_mappings から取得）
- トピックマップ（node_aliases の canonical ノード一覧）
- Embedding 類似度上位ノード（処理対象ノード・チャンクに近いもの上位 N 件）

Layer 3 は初期から注入する。関連性が低いコンテキストはノイズになるため、Embedding 類似度でフィルタリングしてから渡す。

---

## Retry Policy

- Gemini の返却 JSON が不正な場合は JSON repair を 1 回だけ試行する
- JSON repair 後も不正な場合、同一入力に対する Gemini 再試行を最大 2 回まで行う
- LLM 呼び出し自体が失敗した場合も Gemini 再試行を最大 2 回まで行う
- JSON repair は syntax error のみを対象とし、semantic error は補正しない
- semantic error とは、JSON としては読めるが schema・enum・参照整合性・level 制約を満たさない状態を指す
- `documents.status` を `failed` に更新し、失敗理由をログに記録する
- 再処理は `StartProcessing` の `force_reprocess=true` で対応する
- 評価データ（[13-evaluation-data.md](13-evaluation-data.md)）を使った品質劣化検知で根本原因を特定する

### JSON Repair の対象範囲

- 許容する repair:
  - Markdown code fence の除去
  - 末尾カンマの除去
  - 閉じ括弧・閉じ角括弧の不足補完
  - クォート崩れなどの軽微な JSON 構文修正
- 許容しない repair:
  - `level=8` を `3` に補正する
  - `category=\"foo\"` を既知 enum に寄せる
  - 不明な `edge_type` を `related_to` に寄せる
  - 欠落した `label` や `source` を推測補完する

### フォールバック方針

- 構造成立に必須な項目はフォールバックせず、要素破棄または再試行に回す
- 品質補助項目は保存時にフォールバックを許容する
- 初期実装で許容する主なフォールバックは以下とする:
  - `description` 欠落 → 空文字で保存
  - `summary_html` 欠落 → `null` 保存、フロントは `description` にフォールバック
  - `entity_type` 欠落（`category=entity`）→ `unspecified` で保存

---

## Node Integration Policy（Document 間）

document 内の重複統合は Pass 2 が担う。document 間の統合は以下の順で行う。

1. **Pass 2**: 文書内統合は Gemini が担当し、重複ノードの統合と canonical label の決定を行う
2. **ラベル正規化**: document 間比較では全角/半角、前後空白、大小文字差を吸収する
3. **Embedding 類似度**: Vertex AI Embeddings で `label + description` をベクトル化し、同名異概念の誤統合を避ける
4. **自動承認条件**: 正規化後ラベルの Levenshtein 距離が 2 以下、かつ cosine similarity が 0.97 以上のペアは `node_aliases.merge_status=approved` として自動登録する
5. **要レビュー候補**: cosine similarity が 0.88 以上 0.97 未満のペアは `node_aliases.merge_status=suggested` として登録する
6. **却下条件**: cosine similarity が 0.88 未満のペアは候補を保存しない
7. **承認フロー**: `suggested` は `dev` ロールがレビューし、承認後に canonical_node_id へ読み替え対象とする

embedding は Pass 2 完了直後に生成し、document 処理パイプラインの一部として保存する。

---

### フィールド

| フィールド | 型 | 説明 |
| --- | --- | --- |
| `level` | INT64 | 0〜3（ドメイン→概念→施策→詳細） |
| `category` | STRING | `concept` / `entity` / `claim` / `evidence` / `counter` |
| `entity_type` | STRING | `organization` / `person` / `metric` / `date` / `location`（category=entity のみ） |

---

## Edge Type System

### EdgeType

| edge_type | 説明 |
| --- | --- |
| `hierarchical` | 階層の親子関係 |
| `supports` | evidence が claim を支持する |
| `contradicts` | counter が claim に反論する |
| `related_to` | 汎用的な関連（分類できない場合） |
| `measured_by` | concept/entity が metric で測定される |
| `involves` | concept に entity が関与する |
| `causes` | concept が別の concept を引き起こす |
| `exemplifies` | 上位概念の具体例 |

---

## LLM 呼び出し数の見積もり

1ドキュメントあたりの概算：

| ステージ | 呼び出し数 |
| --- | --- |
| セマンティックチャンキング | 1回 |
| Pass 1（チャンク数に比例） | チャンク数 × 1回（例: 10チャンク → 10回） |
| Pass 2（文書全体統合） | 1回 |
| HTML サマリ生成（ノード数に比例） | ノード数 × 1回（例: 30ノード → 30回） |
| **合計** | **約 42回 / ドキュメント（10チャンク・30ノードの場合）** |

HTML サマリはノード数が多い場合にコストが支配的になるため、バッチ化または並列化を検討する。

---

## Open Issues

- Pass 2 のコンテキスト長上限（大きな文書では Pass 1 全結果が入りきらない可能性）
- HTML サマリ生成の並列化戦略（Cloud Tasks でノード単位に並列投入するか）
- level 割り当ての一貫性確保（文書間で同じ概念が異なる level になるケースの対処）
