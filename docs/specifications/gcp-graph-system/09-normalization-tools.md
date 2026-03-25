# 09. Normalization Tools

## Purpose

正規化ツール層は、文書のノード化前に構造や型の乱れを補正するための再利用可能な処理単位を提供する。主用途は、LLM が Python スクリプト案を生成し、それをレビュー、保存、再実行可能なツールとして管理することである。

## Target Use Cases

- 日本語データで数値が文字列として保存されている場合の型補正
- 全角数字、全角記号、表記揺れの正規化
- 壊れた CSV や JSON の列補正
- 日付形式や列名の統一
- PDF 抽出後のノイズ除去

## Tool Lifecycle

1. 問題データまたは問題説明を入力する
2. LLM が Python スクリプト案を生成する
3. サンドボックスで dry-run する
4. 差分とログを確認する
5. 保存して version を付与する
6. 承認済みツールとして再利用する
7. 必要に応じて改訂、廃止する

## Tool Package Layout

```text
tools/
└─ normalize_jp_numeric_strings/
   ├─ tool.py
   ├─ manifest.yaml
   ├─ README.md
   └─ fixtures/
      ├─ input.json
      └─ expected.json
```

## Manifest Requirements

- `tool_id`
- `name`
- `version`
- `description`
- `input_format`
- `output_format`
- `allowed_file_types`
- `timeout_sec`
- `memory_limit_mb`
- `created_by`
- `approval_status`

## Execution Model

- ツールは Python として実装する
- 実行はサンドボックスコンテナ内で行う
- ネットワークアクセスは原則禁止とする
- 読み取りと書き込みの対象ディレクトリを限定する
- dry-run と本実行を明確に分離する

## Approval States

- `draft`
- `reviewed`
- `approved`
- `deprecated`

## Outputs

- 正規化済み成果物
- 標準出力、標準エラー
- 変換差分
- 実行時間
- exit status

## Risks

- LLM が危険な Python を生成する可能性
- 差分が大きすぎる変換による意図しない破壊
- 同一ツールの version 管理不足

## Mitigations

- サンドボックス実行
- approval 状態を必須にする
- 原本は必ず保持する
- dry-run を標準フローにする
- 変換前後差分を保存する
