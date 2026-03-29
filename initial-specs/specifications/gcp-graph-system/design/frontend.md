# 15. Frontend Specification

## Overview

フロントエンドは `TypeScript + React + React Flow` で実装し、`Firebase Hosting` から配信する。グラフの可視化と対話的探索を中心に、ファイルアップロード・処理ステータス確認・ノード詳細閲覧・近傍展開・多段経路検索のユースケースを担う。開発者向けに `/dev/stats` ルートで統計ビューワーを提供し、`Firebase Auth` のカスタムクレーム `role: "dev"` を持つユーザーのみ表示する。

## UI State Naming Policy

フロントの状態名は [data-model.md](../domain/data-model.md) の state family を踏まえつつ、UI 固有の family を追加して扱う。

### Frontend State Families

| Family | 用途 | 主な値 |
| --- | --- | --- |
| `DocumentLifecycleState` | document 一覧・再処理ボタンの表示制御 | `uploaded` / `pending_normalization` / `processing` / `completed` / `failed` |
| `GraphProjectionScope` | node / edge が document graph か canonical graph かの区別 | `document` / `canonical` |
| `GraphViewMode` | canvas のレイアウト表示モード | `level_bands` / `force` / `claim_focus` |
| `ExplorePanelState` | 詳細パネルの表示状態 | `closed` / `document_detail` / `canonical_detail` |
| `PathSearchMode` | 経路検索 UI の操作状態 | `inactive` / `picking_source` / `picking_target` / `results` |
| `ExploreDepthPreset` | 近傍展開の深さプリセット | `one_hop` / `two_hop` / `three_hop` |

### Naming Rules

- backend 由来の `status` は family 名付きで扱い、`documentStatus` のように用途を明示する
- `scope` は `GraphProjectionScope` に限定し、単なる UI タブ切り替えには使わない
- `mode` は `GraphViewMode` や `PathSearchMode` のように具体 family 名で分ける
- React state 名は family を反映し、`viewMode`, `pathSearchMode`, `explorePanelState` のように曖昧な `state` 単体名を避ける

---

## ルート構成

| パス | 用途 |
| --- | --- |
| `/` | ドキュメント一覧・グラフビュー |
| `/dev/stats` | 開発者向け統計ビューワー（`role: "dev"` のみ表示） |

---

## 画面構成

```
┌─────────────────────────────────────────────────┐
│ Header: ActionRev      [ドキュメント一覧] [経路検索]│
├──────────────┬──────────────────────────────────┤
│              │ [レベル帯▼] [力学] [Claim] [1-hop] 🔍 │
│  ドキュメント │ ─────────────────────────────── │
│  一覧         │                                  │
│              │         Graph Canvas             │
│  ・doc_001  │                                  │
│  ・doc_002  │                                  │
│  ・doc_003  │                                  │
│              │                                  │
│  [+ アップ   │                                  │
│    ロード]   │                                  │
├──────────────┴──────────────────────────────────┤
│ Node Detail / Explore Panel（ノードクリック時に展開）│
└─────────────────────────────────────────────────┘
```

---

## グラフレイアウト

### デフォルト: レベル帯レイアウト

水平帯でレベルを固定し、横方向は力学で整列する。

```
━━━━━━━━━━━━━━━━━ level 0 ━━━━━━━━━━━━━━━━━━━
              [事業戦略]

━━━━━━━━━━━━━━━━━ level 1 ━━━━━━━━━━━━━━━━━━━
      [販売戦略]              [マーケ戦略]

━━━━━━━━━━━━━━━━━ level 2 ━━━━━━━━━━━━━━━━━━━
  [テレアポ]  [既存深耕]          [SNS施策]

━━━━━━━━━━━━━━━━━ level 3 ━━━━━━━━━━━━━━━━━━━
 [スクリプト改善]  [CV率3.2%]     [A社事例]
```

- `hierarchical` エッジは縦方向
- `related_to` / `causes` などの横断エッジは斜め・水平で描画
- 実装ライブラリ: `elkjs`（ELK layered アルゴリズム）

### ビュー2: 力学モデル

複数ドキュメント横断の探索に使う。関連するノードが自然にクラスター化する。

- 実装ライブラリ: `d3-force` + React Flow カスタムレイアウト
- ノード fill は `category` 色を維持し、border color を `document_id` ごとに割り当てる
- `source_filename` / `document_id` でノードをグループ化するオプションを用意する

### ビュー3: Claim 中心ビュー

`claim` ノードを中心に `evidence` / `counter` を展開して表示する。

```
         [claim: SNSの方がROIが高い]
        ↙ supports      ↘ contradicts
[A社 CV率3.2%]     [テレアポは関係構築に強み]
```

- `claim` / `evidence` / `counter` category のノードのみ表示する
- 実装: dagre でクレームを中心ノードとして配置

### 切り替え

Canvas 右上のボタンでビューを切り替える。

```
[レベル帯 ▼]  [力学]  [Claim]
```

探索系の操作は同じツールバーに配置する。

```
[1-hop] [2-hop] [3-hop] [経路検索]
```

- `1-hop` / `2-hop` / `3-hop` は選択中ノードに対する近傍展開の深さ
- `経路検索` は 2 ノード選択モードへ切り替える
- この切り替えは `GraphViewMode` と `PathSearchMode` の変更として扱う

---

## ノードの見た目

### サイズ（level に対応）

| level | サイズ |
| --- | --- |
| 0 | 180px × 60px |
| 1 | 140px × 50px |
| 2 | 110px × 40px |
| 3 | 80px × 32px |

### 色（category に対応）

| category | 色 |
| --- | --- |
| `concept` | 青系 `#4A90D9` |
| `entity` | 緑系 `#5BAD6F` |
| `claim` | オレンジ系 `#E8933A` |
| `evidence` | 黄系 `#D4B84A` |
| `counter` | 赤系 `#D95A5A` |

### ラベル表示

- ノード内に `label` を表示する
- `level=3` はテキストを小さくする（10px）
- `summary_html` が null のノードはラベルをイタリック表示してフォールバックを示す
- `scope=canonical` のノードは二重 border と `canonical` バッジで document ノードと区別する

---

## エッジのスタイル

| edge_type | 線種 | 色 |
| --- | --- | --- |
| `hierarchical` | 実線 | グレー `#888` |
| `supports` | 点線 | 黄緑 `#7CB87C` |
| `contradicts` | 波線 | 赤 `#D95A5A` |
| `related_to` | 細点線 | ライトグレー `#BBB` |
| `measured_by` | 実線（細） | 緑 `#5BAD6F` |
| `involves` | 実線（細） | 青 `#4A90D9` |
| `causes` | 矢印付き実線 | 紫 `#9B6DD9` |
| `exemplifies` | 破線 | グレー `#AAA` |

- `scope=canonical` の edge は glow またはやや太い線で描画し、document edge と区別する

---

## ノード詳細・探索パネル

ノードをクリックすると画面下部にパネルが展開する。

```
┌─────────────────────────────────────────────────┐
│ [concept / level 2] テレアポ施策   [1-hop][2-hop][起点]│
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌── iframe (summary_html) ──────────────────┐  │
│  │ <h3>概要</h3>                              │  │
│  │ <table>...</table>                         │  │
│  │ <ul>...</ul>                               │  │
│  └────────────────────────────────────────────┘  │
│                                                  │
│  出典チャンク                                     │
│  ┌────────────────────────────────────────────┐  │
│  │ [report.pdf / p.3]                         │  │
│  │ "テレアポ施策として月次100件を目標に..."    │  │
│  └────────────────────────────────────────────┘  │
│                                                  │
│  探索操作                                         │
│  [近傍を展開] [このノードを経路検索の起点にする]   │
│  [このノードを終点にする] [展開を折りたたむ]       │
│                                                  │
│  関連ノード                                       │
│  [販売戦略 ↑]  [スクリプト改善 ↓]  [CV率3.2% →] │
└─────────────────────────────────────────────────┘
```

- `summary_html` は `<iframe sandbox="allow-same-origin" srcdoc={summary_html}>` で描画する
- CSS は `iframe` の `load` 後に `iframe.contentDocument` へ `<link rel="stylesheet">` を動的注入する
- `summary_html` が null の場合は `description` をプレーンテキストで表示する
- 出典チャンクには `source_filename` とページ番号を表示する
- 関連ノードは隣接エッジを辿って取得し、クリックで該当ノードにフォーカスする
- 探索操作の `近傍を展開` は `ExpandNeighbors` を呼び、取得した subgraph を現在の canvas に追加する
- `起点` / `終点` を設定後に `経路検索` を実行すると `FindPaths` を呼び、結果の path をハイライト表示する
- 展開済み subgraph は document の初期 graph と区別できるよう外枠または薄い背景で表示する
- `scope=canonical` のノード詳細では `canonical_node_id` と代表ラベルを表示し、document ノード詳細では `document_id` と出典 chunk を優先表示する
- `scope=canonical` の edge を含む path は優先ハイライトし、document edge は補助的に薄く残す
- 詳細パネルの取得 API は `GetGraphEntityDetail` に統一し、`target_ref.scope` と `target_ref.id` を切り替えて document / canonical の両方を扱う
- canonical 詳細では `representative_nodes` と `evidence.supporting_edges` を表示し、path の各関係から document 根拠へ降りられるようにする

---

## 対話的探索 UX

### 近傍展開

- ユーザーは任意のノードをクリックして `1-hop` / `2-hop` / `3-hop` 展開を実行できる
- 展開結果は既存 graph に追加マージし、既存ノードと重複するノードは再利用する
- 展開したノード・エッジはアニメーション付きでキャンバスに出現する
- 追加分は `edge_type` ごとにフィルタできる
- 展開の undo / collapse をサポートする

### 経路検索

- ユーザーは 2 ノードを `起点` / `終点` として選択できる
- `経路検索` 実行後、候補 path をサイドパネルまたは上部トレイに一覧表示する
- path を選ぶと対応ノード・エッジのみ強調表示し、他は半透明にする
- 表示内容には hop 数、edge type の列、document 横断かどうかを含める
- path 一覧には `source_document_ids` を併記し、`supporting_edge_ids` がある path は「根拠あり」として表示する

### document 横断探索

- 初期 graph は document 単位でロードする
- 近傍展開・経路検索では canonical node を起点に document 横断の subgraph を追加表示できる
- document 横断で取得したノードは border color で document ごとの差異を示す
- 現在表示中の graph に対し `初期 document のみ` / `横断を含む` のトグルを用意する

---

## フロント状態設計

探索 UX のため、フロントは「初期 graph」と「探索で追加した graph」を分けて保持する。

### state の分離

- `baseGraph`: `GetGraph` で取得した document 初期表示用 graph
- `expandedGraph`: `ExpandNeighbors` で追加した subgraph
- `pathSearchGraph`: `FindPaths` の結果として一時追加した subgraph
- `selectedNodeId`: 現在詳細表示しているノード
- `pathSearchDraft`: 経路検索の `source_node_id` / `target_node_id`
- `exploreOptions`: `maxDepth`, `edgeTypeFilters`, `crossDocument`
- `viewMode`: `GraphViewMode`
- `explorePanelState`: `ExplorePanelState`
- `pathSearchMode`: `PathSearchMode`
- `exploreDepthPreset`: `ExploreDepthPreset`

### state 更新ルール

- `GetGraph` 実行時は `baseGraph` を置き換え、`expandedGraph` / `pathSearchGraph` をクリアする
- `ExpandNeighbors` 実行時は `expandedGraph` にマージする
- `FindPaths` 実行時は `pathSearchGraph` を置き換える
- 同一ノード ID の重複はフロント側で統合し、最初に読み込んだ node オブジェクトを基準に不足属性だけ補完する
- `scope=document` と `scope=canonical` は別ノードとして保持し、`canonical_node_id` を使って関連表示だけを行う
- `edge.scope=document` と `edge.scope=canonical` も別系列として保持し、path overlay では `scope=canonical` を優先する
- `collapse` 実行時は対象操作で追加した subgraph のみを取り除く
- ノード選択時は `explorePanelState` を `document_detail` または `canonical_detail` に遷移させる
- `経路検索` 押下時は `pathSearchMode` を `picking_source` にし、起点選択後に `picking_target`、結果取得後に `results` に遷移させる

### 表示上のレイヤー

- `baseGraph`: 通常表示
- `expandedGraph`: 追加表示。薄いハイライト背景または発光 border を付ける
- `pathSearchGraph`: 最上位表示。path 上のノードと edge を強調し、非該当要素は減衰表示する

---

## ファイルアップロード UI

```
┌─────────────────────────────────────────┐
│                                          │
│   ドラッグ＆ドロップ または クリック    │
│                                          │
│   対応: PDF / Markdown / TXT / CSV / zip │
│                                          │
└─────────────────────────────────────────┘
```

- アップロード後、即時処理開始
- zip の場合は処理中に展開されたファイル一覧を表示する
- `extraction_depth` を選択できるトグルを用意する（デフォルト: `full`）

---

## 処理ステータス表示

ドキュメント一覧に処理状態をバッジで表示する。

- この表示は `DocumentLifecycleState` に直接対応する

| status | 表示 |
| --- | --- |
| `uploaded` | ⬆ アップロード済 |
| `pending_normalization` | 🛠 承認待ち |
| `processing` | ⏳ 処理中（スピナー） |
| `completed` | ✅ 完了 |
| `failed` | ❌ 失敗（再処理ボタンを表示） |

zip の場合は展開されたファイルごとの処理状況も折りたたみで確認できる。

---

## フィルタ・検索

Canvas 右上の 🔍 アイコンで展開するサイドパネル。

- `node_type_filter`: category で絞り込む
- `source_filename_filter`: zip 内ファイルで絞り込む
- テキスト検索: label / description の部分一致
- `edge_type_filter`: 近傍展開・経路検索に含める edge_type を選ぶ
- `depth_filter`: 近傍展開の hop 数を選ぶ
- `cross_document_toggle`: document 横断ノードを表示するか切り替える
- `document_scope_filter`: 探索対象 document を現在文書のみ / 選択文書群 / workspace 全体 から選ぶ

---

## レイアウト切り替えアニメーション

ビューを切り替えた際、ノードが現在位置から新しい位置へアニメーションで移動する。

- 新しいレイアウトの座標を計算し、React Flow の `setNodes` で位置を更新する
- React Flow の CSS transition（`transition: transform 300ms ease-in-out`）で補間する
- 切り替え中は操作を無効化し、完了後に有効化する
- ノード数が多い（100件以上）場合はアニメーションをスキップしてパフォーマンスを優先する
- 近傍展開時は追加ノードのみフェードインし、既存ノードの大きな再配置は避ける

---

## 経路検索モード

```
┌─────────────────────────────────────────────┐
│ 経路検索: [起点: 販売戦略] [終点: CV率3.2%]  │
│ [max depth: 4▼] [edge type: すべて▼] [検索]  │
├─────────────────────────────────────────────┤
│ Path 1  販売戦略 → SNS施策 → CV率3.2%  (2 hop) │
│ Path 2  販売戦略 → テレアポ → A社事例 ...      │
└─────────────────────────────────────────────┘
```

- 経路検索モード中は canvas 上でノードを 2 つまで選択できる
- `max depth` は 2〜6 の範囲で選択可能にする
- path 候補をクリックすると対応サブグラフを中央にフィット表示する
- `supporting_edge_ids` がある場合は「根拠を見る」導線を表示し、必要時に `GetGraphEntityDetail` で詳細を引く
- 検索時は `cross_document_toggle` と `document_scope_filter` の両方を request に反映する

---

## 大量ノード時のパフォーマンス対策

ズームレベルに連動した LOD（Level of Detail）で表示ノード数を制御する。

| ズームレベル | 表示する level |
| --- | --- |
| < 0.4 | level 0 のみ |
| 0.4 〜 0.7 | level 0〜1 |
| 0.7 〜 1.0 | level 0〜2 |
| > 1.0 | 全レベル（level 0〜3） |

- 非表示ノードへのエッジは最近接の表示ノードに折りたたんで描画する
- React Flow の `nodeTypes` でカスタムノードを使い、zoom イベントで表示切り替えを行う
- 500件を超える場合は同一 level・category のノードをクラスターノードに集約し、クリックで展開できるようにする

---

---

## 開発者向け統計ビューワー（`/dev/stats`）

### 画面構成

```
┌──────────────────────────────────────────────────────────────────┐
│ /dev/stats   [パイプライン][抽出品質][評価][エラー][正規化ツール] │
├──────────────────────────────────────────────────────────────────┤
│                       タブ切り替え                                │
└──────────────────────────────────────────────────────────────────┘
```

---

### タブ 1: パイプライン統計

```
┌─── ステージ別 処理時間 ────────────────────────────┐
│  semantic chunking   ████░░░░  avg 1.2s  p95 3.1s  │
│  pass 1 / chunk      ████████  avg 2.4s  p95 5.8s  │
│  pass 2              ██████░░  avg 3.1s  p95 7.2s  │
│  HTML summary        █████░░░  avg 1.8s  p95 4.3s  │
└────────────────────────────────────────────────────┘

┌─── 処理結果 ─────┐   ┌─── Gemini 呼び出し/コスト（日別）─┐
│ completed  94%   │   │  ▁▃▅▇▅▃▇▅  呼び出し数            │
│ failed      6%   │   │  ▁▂▃▄▃▂▄▃  コスト ($)            │
│ → 失敗一覧へ     │   └──────────────────────────────────┘
└──────────────────┘
```

**データソース**
- `documents.status` / `documents.created_at` / `documents.updated_at`
- `processing_jobs`（ステージ別の開始・完了時刻）
- Gemini 呼び出しログ（Cloud Logging から集計）

---

### タブ 2: 抽出品質統計

```
┌─── ドキュメント選択 ──────────────────────────────┐
│  [全体 ▼]  または  [doc_001 ▼]                   │
└──────────────────────────────────────────────────┘

┌─── level 分布 ────────────┐  ┌─── category 分布 ─────────┐
│  level 0  ██░░░░░   3件   │  │  concept  ████████  58%   │
│  level 1  ████░░░  12件   │  │  entity   █████░░░  21%   │
│  level 2  ████████ 34件   │  │  claim    ███░░░░░  12%   │
│  level 3  ████████ 87件   │  │  evidence ██░░░░░░   7%   │
└───────────────────────────┘  │  counter  █░░░░░░░   2%   │
                                └───────────────────────────┘

┌─── Pass 1 → Pass 2 統合数 ──────────────────────┐
│  Pass 1 抽出: 136ノード                           │
│  Pass 2 後:   122ノード  (-14件 統合)             │
│  エッジ補完:  +8件                               │
└──────────────────────────────────────────────────┘
```

**データソース**
- `nodes` テーブル（level / category 集計）
- `processing_jobs`（Pass 1 / Pass 2 のノード数差分）

---

### タブ 3: 評価トレンド

```
┌─── Precision / Recall 週次推移 ─────────────────┐
│  1.0 │                                           │
│  0.8 │  ·─·─·  Precision                        │
│  0.6 │  ○─○─○  Recall                           │
│  0.4 │                                           │
│      └─── w1 ─── w2 ─── w3 ─── w4              │
└──────────────────────────────────────────────────┘

┌─── level 割り当て一致率 ─────────────────────────┐
│  level 0  ████████████████████  95%              │
│  level 1  ████████████████░░░░  88%              │
│  level 2  ██████████████░░░░░░  71%              │
│  level 3  ████████████░░░░░░░░  63%              │
└──────────────────────────────────────────────────┘

┌─── プロンプト変更ログ ───────────────────────────┐
│  2026-03-20  コンテキスト注入 Layer 3 追加  +3% P │
│  2026-03-14  Pass 2 プロンプト改訂         +5% R │
└──────────────────────────────────────────────────┘
```

**データソース**
- `eval/` の gold document 自動評価結果（BigQuery に記録）
- プロンプト変更ログ（git commit と紐付け）

---

### タブ 4: エラー分析

```
┌─── エラー種別 ────────────┐  ┌─── 失敗ドキュメント一覧 ──────────────┐
│  JSON parse    42%  ████  │  │  doc_023  JSON parse error  03-27  [再処理]│
│  timeout       31%  ███   │  │  doc_031  Gemini timeout    03-26  [再処理]│
│  BQ write      27%  ██    │  │  doc_019  BQ write error    03-24  [再処理]│
└───────────────────────────┘  └────────────────────────────────────────────┘

┌─── エラーメッセージ詳細（行クリックで展開）───────────────────────────────┐
│  doc_023 / 2026-03-27 14:32                                              │
│  stage: pass_1  chunk_index: 3                                           │
│  error: unexpected token at position 142: ...                            │
└──────────────────────────────────────────────────────────────────────────┘
```

**データソース**
- `documents.status=failed` + Cloud Logging のエラーログ
- `[再処理]` ボタンは `StartProcessing(force_reprocess=true)` を呼び出す

---

---

### タブ 5: 正規化ツール管理

`dev` ロールのみ表示。workspace とは無関係のシステムグローバルなツール一覧。

```
┌─── 正規化ツール一覧 ──────────────────────────────────────────────┐
│  名前                       状態       問題パターン         操作   │
│  normalize_mojibake_shiftjis  approved  Shift-JIS文字化け   [実行] │
│  fix_csv_columns              reviewed  CSV列ずれ      [承認][実行] │
│  remove_pdf_noise             draft     PDFノイズ  [dry-run][削除] │
│                                                    [+ 新規生成]    │
└───────────────────────────────────────────────────────────────────┘
```

**ツール詳細・dry-run 結果画面：**
```
┌─── normalize_mojibake_shiftjis ────────────────────────────────────┐
│ 問題パターン: Shift-JIS で作成されたファイルを UTF-8 として読んだ際の文字化け │
│ 状態: reviewed    作成者: dev@example.com                           │
│                                                                     │
│ [dry-run 結果]                                                      │
│  変換前                  変換後                                     │
│  "鐚懿・邏・"      →    "日本語テキスト"                           │
│  変更行数: 47 / 312 行                                              │
│                                                                     │
│               [承認する]  [廃止にする]  [再 dry-run]               │
└─────────────────────────────────────────────────────────────────────┘
```

状態遷移ボタン：
- `draft` → `[dry-run]` のみ表示
- `reviewed` → `[承認する]` `[再 dry-run]` を表示
- `approved` → `[本実行]` `[廃止にする]` を表示

---

### API 追加

`GraphService` / `DocumentService` に以下を追加する。

| RPC | 用途 |
| --- | --- |
| `ExpandNeighbors` | 選択ノードから指定 hop の近傍 subgraph を取得 |
| `FindPaths` | 2 ノード間の多段経路を取得 |
| `GetPipelineStats` | ステージ別処理時間・成功率・Gemini コスト集計 |
| `GetExtractionStats` | ノード level/category 分布・Pass 統合数 |
| `GetEvaluationTrend` | 週次 Precision/Recall・level一致率 |
| `ListFailedDocuments` | 失敗ドキュメント一覧とエラー詳細 |

### フロント実装メモ

- 初期ロードでは `GetGraph` のみを呼び、node ごとの追加 fetch は行わない
- ノード詳細パネルを開いた時のみ `GetGraphEntityDetail` を呼ぶ
- 近傍展開は `ExpandNeighbors`、経路検索は `FindPaths` を明示的なユーザー操作でのみ呼ぶ
- 探索中も `GetGraph` を再実行しない限り base graph は保持する

---

## Open Issues

- モバイル対応は初期スコープ外とするが、将来 React Flow の簡易閲覧 UI を別設計にするか
