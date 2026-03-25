# 08. AI Pipeline

## Overview

AI パイプラインは、原本の保存からグラフ保存までを段階的に処理する。初期段階では単一 document 処理を前提とし、後続で非同期ジョブ化と並列化を行う。

## Stages

### 1. Raw Intake

- 元ファイルを `Cloud Storage` に保存する
- 原本は上書きせず、再処理可能な状態を維持する

### 2. Normalization

- データ不整合がある場合、LLM が Python 正規化スクリプト案を生成する
- 既存ツールがある場合は再利用を優先する
- スクリプトはサンドボックスで dry-run し、差分を確認可能にする
- 承認済みツールのみ本実行できるようにする
- 正規化済み成果物を別保存する

### 3. Extraction

- 正規化済み document からテキストを抽出する
- 文書を chunk に分割する
- chunk ごと、または chunk 群ごとに Gemini に投入する

### 4. Structuring

- Gemini の JSON 出力を受け取る
- ラベル、型、不正 JSON を正規化する
- 重複ノードを統合する
- 出典 chunk を各 node、edge に関連付ける

### 5. Persistence

- `documents`
- `document_chunks`
- `nodes`
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
