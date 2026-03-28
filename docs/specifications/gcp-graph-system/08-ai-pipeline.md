# 08. AI Pipeline

## Overview

AI パイプラインは、原本の保存からグラフ保存までを段階的に処理する。初期段階では単一 document 処理を前提とし、後続で非同期ジョブ化と並列化を行う。

## Stages

### 1. Raw Intake

- 元ファイルを `Cloud Storage` に保存する
- アップロード完了後、処理を即時開始する
- zip ファイルの場合は展開し、対応フォーマット（PDF / Markdown / TXT / CSV）のファイルのみを処理対象とする
- zip 内の各ファイルは同一 `document_id` に属する chunk として扱い、`source_filename` で識別する
- 原本は上書きせず、再処理可能な状態を維持する

### 2. Normalization

- データ不整合がある場合、LLM が Python 正規化スクリプト案を生成する
- 既存ツールがある場合は再利用を優先する
- スクリプトはサンドボックスで dry-run し、差分を確認可能にする
- 承認済みツールのみ本実行できるようにする
- 正規化済み成果物を別保存する

### 3. Semantic Chunking

詳細は [12-extraction-strategy.md](12-extraction-strategy.md) を参照。

- 正規化済み document からテキストを抽出する
- 固定サイズ分割ではなく、Gemini が意味の切れ目でチャンクを決定する
- 1チャンクは「1つのトピックまたは論点を扱う単位」とする
- チャンクサイズの上限を超える場合はさらに分割する
- 各チャンクに `heading`（セクション見出し相当）を付与する

### 4. Pass 1 — Fine-grained Extraction（chunk 単位）

詳細は [12-extraction-strategy.md](12-extraction-strategy.md) を参照。

- 各 chunk に対して個別に Gemini を呼び出す
- `concept` / `entity` / `claim` / `evidence` / `counter` を抽出する
- エッジ種別（`supports`, `contradicts`, `measured_by`, `involves`, `causes`, `exemplifies`）を付与する
- 出典 chunk を各 node / edge に関連付ける

### 5. Pass 2 — Document-level Synthesis（文書全体統合）

詳細は [12-extraction-strategy.md](12-extraction-strategy.md) を参照。

- Pass 1 の全チャンク抽出結果をまとめて Gemini に投入する
- 重複ノードを統合し、canonical ラベルを決定する
- 各ノードに level（0〜3）を付与する
- chunk をまたぐ関係（`hierarchical`, `causes` など）を補完する
- クレーム構造（claim / evidence / counter の論理関係）を整理する

### 6. HTML Summary Generation

- ノード抽出・統合が完了した後、別の Gemini 呼び出しでノードごとの HTML サマリを生成する
- Gemini には `<table>`, `<ul>`, `<h3>` などの構造タグのみを出力させ、`style` 属性・`<style>` タグは含めない
- 生成した HTML は `nodes.summary_html` に保存する
- 生成失敗時は null のまま保存し、フロントは `description` にフォールバックする
- フロントは `<iframe sandbox="allow-same-origin" srcdoc={summary_html}>` で描画し、CSS はアプリ側から注入する

### 7. Persistence

- `documents`
- `document_chunks`
- `nodes`（`summary_html` を含む）
- `edges`
- 将来的には `processing_jobs`, `normalization_tools`, `normalization_tool_runs`

## Design Principles

- 原本は不変とする
- LLM には直接データ変換をさせず、可能な限り再利用可能なツールを生成させる
- 変換処理は dry-run と本実行を分離する
- 差分、ログ、失敗理由を追跡可能にする
- ノード化より前に正規化層を置く

## Future Enhancements

- chunk 並列化
- 正規化ツールの自動候補提示
- ツール選択の類似ケース推薦
- 正規化ルールからの半自動テスト生成
